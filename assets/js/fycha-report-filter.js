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
})();
