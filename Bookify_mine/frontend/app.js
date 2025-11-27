const form = document.getElementById('query-form');
const baseUrlInput = document.getElementById('base-url');
const bookIdInput = document.getElementById('book-id');
const userIdInput = document.getElementById('user-id');
const statusEl = document.getElementById('status');
const resultsSection = document.getElementById('results-section');
const bookResult = document.getElementById('book-result');
const playlistResult = document.getElementById('playlist-result');
const reviewsResult = document.getElementById('reviews-result');
const shelfResult = document.getElementById('shelf-result');
const rawJson = document.getElementById('raw-json');

const STORAGE_KEY = 'bookify.gateway.url';
const DEFAULT_BASE_URL = 'http://localhost:8080';
const shelfLabels = {
  0: 'Sin shelf asignado',
  1: 'Por leer',
  2: 'Leyendo',
  3: 'Leído',
};

const savedBaseUrl = localStorage.getItem(STORAGE_KEY);
baseUrlInput.value = savedBaseUrl || DEFAULT_BASE_URL;

baseUrlInput.addEventListener('blur', () => {
  localStorage.setItem(STORAGE_KEY, baseUrlInput.value.trim());
});

resetPanels();

form.addEventListener('submit', async (event) => {
  event.preventDefault();
  const baseUrl = sanitizeBaseUrl(baseUrlInput.value.trim());
  const bookId = bookIdInput.value.trim();
  const userId = userIdInput.value.trim();

  if (!bookId) {
    setStatus('Ingresa el ID del libro para continuar.', 'error');
    bookIdInput.focus();
    return;
  }

  setStatus('Consultando microservicios…', 'info');
  toggleFormDisabled(true);

  try {
    const data = await fetchOverview(baseUrl, bookId, userId);
    renderData(data);
    setStatus('Datos sincronizados correctamente.', 'success');
  } catch (error) {
    console.error(error);
    setStatus(error.message || 'Algo salió mal, revisa los logs.', 'error');
  } finally {
    toggleFormDisabled(false);
  }
});

function resetPanels() {
  const emptyMessage = '<p class="muted">Aún sin datos disponibles.</p>';
  bookResult.innerHTML = emptyMessage;
  playlistResult.innerHTML = emptyMessage;
  reviewsResult.innerHTML = emptyMessage;
  shelfResult.innerHTML = emptyMessage;
  rawJson.textContent = '—';
}

function toggleFormDisabled(disabled) {
  Array.from(form.elements).forEach((el) => {
    if (el instanceof HTMLButtonElement || el instanceof HTMLInputElement) {
      el.disabled = disabled;
    }
  });
}

async function fetchOverview(baseUrl, bookId, userId) {
  const params = new URLSearchParams({ bookId });
  if (userId) {
    params.append('userId', userId);
  }
  const url = `${baseUrl}/overview/book?${params.toString()}`;
  const response = await fetch(url, { headers: { Accept: 'application/json' } });
  if (!response.ok) {
    const message = await response.text();
    throw new Error(`Gateway respondió ${response.status}: ${message}`);
  }

  return response.json();
}

function renderData(payload) {
  renderBook(payload.book);
  renderPlaylist(payload.playlist);
  renderReviews(payload.reviews);
  renderShelf(payload.shelf);
  rawJson.textContent = JSON.stringify(payload, null, 2);
}

function renderBook(book) {
  if (!book) {
    bookResult.innerHTML = '<p class="muted">Sin resultados del servicio library.</p>';
    return;
  }
  const rows = [
    ['ID', book.id],
    ['Título', book.title],
    ['Autor', book.author],
    ['Páginas', book.pages],
    ['Edición', book.edition],
  ];
  bookResult.innerHTML = `
    <ul class="detail-list">
      ${rows
        .map(
          ([label, value]) => `
            <li>
              <span>${escapeHtml(label)}</span>
              <span>${escapeHtml(value ?? '—')}</span>
            </li>
          `,
        )
        .join('')}
    </ul>
  `;
}

function renderPlaylist(playlist) {
  if (!playlist || !Array.isArray(playlist.tracks) || !playlist.tracks.length) {
    playlistResult.innerHTML = '<p class="muted">No hay playlist disponible.</p>';
    return;
  }
  playlistResult.innerHTML = `
    <p class="muted">Tracks para ${escapeHtml(playlist.bookId || 'el libro')}</p>
    <div class="song-list">
      ${playlist.tracks
        .map(
          (track) => `
            <div class="badge info">
              <span>${escapeHtml(track.title || 'Sin título')}</span>
              <span>·</span>
              <span>${escapeHtml(track.artist || 'Artista desconocido')}</span>
            </div>
          `,
        )
        .join('')}
    </div>
  `;
}

function renderReviews(reviewsResponse) {
  const reviews = Array.isArray(reviewsResponse?.reviews)
    ? reviewsResponse.reviews
    : reviewsResponse;

  if (!Array.isArray(reviews) || !reviews.length) {
    reviewsResult.innerHTML = '<p class="muted">Sin reseñas por ahora.</p>';
    return;
  }

  const average = (
    reviews.reduce((sum, review) => sum + Number(review.rating || 0), 0) / reviews.length
  ).toFixed(1);

  reviewsResult.innerHTML = `
    <div class="badge rating">Promedio ${average} ⭐ (${reviews.length})</div>
    <div class="review-list">
      ${reviews
        .map(
          (review) => `
            <div class="review-item">
              <strong>${escapeHtml(review.userId || 'Anon')}</strong>
              <div class="muted">${'⭐'.repeat(Math.max(1, Math.min(5, Number(review.rating) || 0)))}</div>
              <p>${escapeHtml(review.text || 'Sin comentario')}</p>
            </div>
          `,
        )
        .join('')}
    </div>
  `;
}

function renderShelf(shelfItem) {
  if (!shelfItem) {
    shelfResult.innerHTML = '<p class="muted">Solicita con un userId para ver el shelf.</p>';
    return;
  }
  const label = shelfLabels[shelfItem.shelf] || `Shelf #${shelfItem.shelf}`;
  shelfResult.innerHTML = `
    <p>Usuario: <strong>${escapeHtml(shelfItem.userId || '—')}</strong></p>
    <p>Libro: <strong>${escapeHtml(shelfItem.bookId || '—')}</strong></p>
    <div class="badge success">${escapeHtml(label)}</div>
  `;
}

function setStatus(message, tone = 'muted') {
  statusEl.textContent = message;
  statusEl.className = `status ${tone}`;
}

function sanitizeBaseUrl(value) {
  if (!value) return DEFAULT_BASE_URL;
  return value.endsWith('/') ? value.slice(0, -1) : value;
}

function escapeHtml(value) {
  return String(value)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;');
}
