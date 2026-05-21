// HTMX global loading bar
document.body.insertAdjacentHTML('afterbegin', '<div id="htmx-indicator"></div>');

// Show/hide loading bar on HTMX requests
document.addEventListener('htmx:beforeRequest', () => {
  document.getElementById('htmx-indicator').style.display = 'block';
});
document.addEventListener('htmx:afterRequest', () => {
  document.getElementById('htmx-indicator').style.display = 'none';
});

// Auto-dismiss alerts after 4s
document.addEventListener('htmx:afterSwap', () => {
  document.querySelectorAll('.alert').forEach(el => {
    setTimeout(() => {
      el.style.transition = 'opacity 0.5s';
      el.style.opacity = '0';
      setTimeout(() => el.remove(), 500);
    }, 4000);
  });
});
