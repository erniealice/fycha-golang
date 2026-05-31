/**
 * fycha — Report filter handler for sheet-based filtering.
 *
 * Handles period preset button clicks and custom date toggling
 * inside the filter sheet (#sheetContent). Uses document-level
 * event delegation since sheet content is loaded dynamically via HTMX.
 */
(function () {
  'use strict';

  // ─── Preset Button Clicks ───
  document.addEventListener('click', function (e) {
    var btn = e.target.closest('.preset-btn');
    if (!btn) return;

    var form = btn.closest('.filter-sheet-form');
    if (!form) return;

    var preset = btn.dataset.preset;
    if (!preset) return;

    // Update active state on all preset buttons in this form
    form.querySelectorAll('.preset-btn').forEach(function (b) {
      b.classList.remove('active');
    });
    btn.classList.add('active');

    // Update hidden period input
    var periodInput = form.querySelector('input[name="period"]');
    if (periodInput) {
      periodInput.value = preset;
    }

    // Toggle custom date inputs
    var customDates = form.querySelector('#custom-dates');
    if (customDates) {
      if (preset === 'custom') {
        customDates.classList.remove('hidden');
      } else {
        customDates.classList.add('hidden');
      }
    }
  });

  // ─── Group-By Radio Styling ───
  document.addEventListener('change', function (e) {
    if (e.target.type !== 'radio' || e.target.name !== 'group-by') return;

    var container = e.target.closest('.groupby-options');
    if (!container) return;

    container.querySelectorAll('.groupby-option').forEach(function (opt) {
      opt.classList.remove('active');
    });
    e.target.closest('.groupby-option').classList.add('active');
  });

  // ─── Filter-Sheet Open (CSP-prep, Plan-6) ───
  // The "Filters" toolbar button (report-filter-btn / report-aging-toolbar-prefix /
  // report-dimension-toolbar-prefix) previously carried inline
  // `onclick="lf.ui.Sheet.open('Filters')"`. Inline event-handler attributes
  // cannot be nonced/hashed and block an enforcing `script-src 'self'`, so the
  // sheet shell is now opened via document-level delegation. The title moved into
  // a `data-sheet-title` attribute (html/template auto-escapes attribute context).
  // The button's own `hx-get`/`hx-target`/`hx-swap` still loads the filter body
  // into #sheetContent independently — same dual behavior as the inline handler.
  // Mirrors the app-level home.js delegation idiom. The `data-sheet-close` path
  // (footer Apply/Clear) is handled globally by pyeza sheet.js — no wiring here.
  if (window.lf && typeof window.lf.on === 'function') {
    lf.on('click', '[data-sheet-open]', function () {
      if (window.lf.ui && window.lf.ui.Sheet && typeof window.lf.ui.Sheet.open === 'function') {
        window.lf.ui.Sheet.open(this.getAttribute('data-sheet-title') || '');
      }
    });
  }
})();
