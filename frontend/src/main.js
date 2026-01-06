import './style.css';

import { EventsOn } from "../wailsjs/runtime/runtime"
import { GetProcesses, GetApplications, HideWindow, FilterItems, SelectItem } from '../wailsjs/go/main/App';

let filteredItems = []
let allItems = []
let selectedIndex = 0
let debounceTimerId = null

const searchInput = document.getElementById("searchInput")
const resultsContainer = document.getElementById("results")
const itemCount = document.getElementById("itemCount")

async function filterItems(query) {
    if (!query) {
        filteredItems = [...allItems];
    } else {
        filteredItems = await FilterItems(query)
    }
    selectedIndex = 0;
    renderResults();
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function debounce(func, delay) {
    return function(...args) {
        clearTimeout(debounceTimerId);
        debounceTimerId = setTimeout(() => func.apply(this, args), delay);
    };
}

const debouncedFilter = debounce(async (query) => {
    await filterItems(query);
}, 200);

function renderResults() {
    if (filteredItems.length === 0) {
        resultsContainer.innerHTML = `<div class="no-results">No results found</div>`
        itemCount.innerText = '0 items'
        return
    }

    resultsContainer.innerHTML = filteredItems.map((item, index) => {
        return `<div class="result-item ${index === selectedIndex ? 'selected' : ''}" 
              data-index="${index}">
            <div class="result-icon">
                <img src="${escapeHtml(item.icon)}" alt="${escapeHtml(item.text)}" />
            </div>
            <div class="result-content">
                <div class="result-text">${escapeHtml(item.text)}</div>
                <div class="result-description">${escapeHtml(item.description)}</div>
            </div>
        </div>`
    }).join('')

    itemCount.textContent = `${filteredItems.length} item${filteredItems.length !== 1 ? 's' : ''}`

    // Scroll selected item into view
    const selected = resultsContainer.querySelector('.result-item.selected')
    if (selected) {
        selected.scrollIntoView({ block: 'nearest', behavior: 'smooth' })
    }
}

async function selectItem() {
    await SelectItem(selectedIndex)
    HideWindow()
}

searchInput.addEventListener("input", (e) => {
    debouncedFilter(e.target.value)
})

searchInput.addEventListener('keydown', (e) => {
    switch (e.key) {
        case 'ArrowDown':
            e.preventDefault();
            selectedIndex = Math.min(selectedIndex + 1, filteredItems.length - 1);
            renderResults();
            break;
        case 'ArrowUp':
            e.preventDefault();
            selectedIndex = Math.max(selectedIndex - 1, 0);
            renderResults();
            break;
        case 'Enter':
            e.preventDefault();
            selectItem();
            break;
        case 'Escape':
            e.preventDefault();
            // Hide window
            HideWindow()
            break;
    }
});

window.addEventListener("load", () => {
    GetApplications()
})

EventsOn("ipc:results", (msg) => {
    allItems = [...msg]
    filteredItems = [...allItems]
    renderResults()
})
