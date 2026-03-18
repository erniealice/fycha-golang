/**
 * Financial Statement — Collapsible Sections + Change% Indicators
 *
 * 1. Collapse: toggles visibility of section body rows on header click.
 * 2. Change%: adds color classes + arrow indicators to .fs-col-change cells.
 *
 * Structure:
 *   <tbody class="fs-section">
 *     <tr class="fs-section-header-row" data-fs-toggle>  <- click target
 *     <tr class="fs-section-body">                        <- hidden/shown
 *     ...
 *   </tbody>
 */
(function () {
    'use strict';

    function initCollapsible() {
        var tables = document.querySelectorAll('.fs-collapsible');

        tables.forEach(function (table) {
            if (table.dataset.fsInit) return;
            table.dataset.fsInit = '1';

            table.addEventListener('click', function (e) {
                var headerRow = e.target.closest('[data-fs-toggle]');
                if (!headerRow) return;

                var section = headerRow.closest('.fs-section');
                if (!section) return;

                section.classList.toggle('fs-collapsed');
            });
        });
    }

    /**
     * Scans all .fs-col-change cells and adds:
     * - fs-change-up   + "▲" prefix for positive values (+X.X%)
     * - fs-change-down + "▼" prefix for negative values (-X.X%)
     * - fs-change-flat for zero (0.0%)
     */
    function initChangeIndicators() {
        var cells = document.querySelectorAll('.fs-col-change');

        cells.forEach(function (cell) {
            // Skip header cells and already-processed cells
            if (cell.tagName === 'TH' || cell.dataset.fsChange) return;

            var text = cell.textContent.trim();
            if (!text) return;

            cell.dataset.fsChange = '1';

            if (text.charAt(0) === '+' && text !== '+0.0%') {
                cell.classList.add('fs-change-up');
                cell.textContent = '\u25B2 ' + text;
            } else if (text.charAt(0) === '-') {
                cell.classList.add('fs-change-down');
                cell.textContent = '\u25BC ' + text;
            } else {
                cell.classList.add('fs-change-flat');
                cell.textContent = '\u2014 ' + text;
            }
        });
    }

    // Init on load
    initCollapsible();
    initChangeIndicators();

    // Re-init after HTMX swaps
    document.addEventListener('htmx:afterSettle', function () {
        initCollapsible();
        initChangeIndicators();
    });
})();
