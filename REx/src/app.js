const DATA_URL = "regex.yaml";

const searchInput = document.getElementById("search-input");
const categorySelect = document.getElementById("category-select");
const clearButton = document.getElementById("clear-button");
const resultsGrid = document.getElementById("results-grid");
const emptyState = document.getElementById("empty-state");
const resultsCount = document.getElementById("results-count");
const statTotal = document.getElementById("stat-total");
const statCategories = document.getElementById("stat-categories");

let records = [];
let categories = [];

const debounce = (fn, delay = 180) => {
  let timer;
  return (...args) => {
    window.clearTimeout(timer);
    timer = window.setTimeout(() => fn(...args), delay);
  };
};

const normalize = (value) => (value || "").toString().toLowerCase();

const buildRecords = (rawGroups) => {
  const collected = [];
  rawGroups.forEach((group) => {
    const groupName = group.name || "Uncategorized";
    const regexes = Array.isArray(group.regexes) ? group.regexes : [];
    regexes.forEach((regexItem) => {
      collected.push({
        id: `${groupName}-${regexItem.name || regexItem.regex}`,
        group: groupName,
        name: regexItem.name || "Untitled",
        regex: regexItem.regex || "",
        example: regexItem.example || "",
        falsePositives: Boolean(regexItem.falsePositives),
        caseinsensitive: Boolean(regexItem.caseinsensitive),
        extraGrep: regexItem.extra_grep || "",
      });
    });
  });
  return collected;
};

const populateCategories = (items) => {
  categories = Array.from(new Set(items.map((item) => item.group))).sort();
  categories.forEach((category) => {
    const option = document.createElement("option");
    option.value = category;
    option.textContent = category;
    categorySelect.appendChild(option);
  });
  statCategories.textContent = categories.length.toString();
};

const renderCard = (item, index) => {
  const card = document.createElement("article");
  card.className = "card";
  card.style.animationDelay = `${Math.min(index * 0.03, 0.3)}s`;

  const title = document.createElement("h3");
  title.textContent = item.name;

  const badgeRow = document.createElement("div");
  badgeRow.className = "badge-row";

  const categoryBadge = document.createElement("span");
  categoryBadge.className = "badge";
  categoryBadge.textContent = item.group;
  badgeRow.appendChild(categoryBadge);

  if (item.falsePositives) {
    const warnBadge = document.createElement("span");
    warnBadge.className = "badge warning";
    warnBadge.textContent = "False positives";
    badgeRow.appendChild(warnBadge);
  }

  if (item.caseinsensitive) {
    const ciBadge = document.createElement("span");
    ciBadge.className = "badge";
    ciBadge.textContent = "Case-insensitive";
    badgeRow.appendChild(ciBadge);
  }

  const pattern = document.createElement("div");
  pattern.className = "pattern";
  pattern.textContent = item.regex;

  const shouldCollapse = item.regex.length > 180 || item.regex.split("\n").length > 6;
  const shouldFullWidth = item.regex.length > 200;
  if (shouldCollapse) {
    pattern.classList.add("collapsed");
  }
  if (shouldFullWidth) {
    card.classList.add("full-width");
  }

  const example = document.createElement("p");
  example.className = "example";
  example.textContent = item.example ? `Example: ${item.example}` : "Example: —";

  const actions = document.createElement("div");
  actions.className = "card-actions";

  const extra = document.createElement("span");
  extra.className = "mono";
  extra.textContent = item.extraGrep ? `extra_grep: ${item.extraGrep}` : "";

  const copyButton = document.createElement("button");
  copyButton.type = "button";
  copyButton.className = "copy-button";
  copyButton.textContent = "Copy regex";
  copyButton.addEventListener("click", async () => {
    try {
      await navigator.clipboard.writeText(item.regex);
      copyButton.textContent = "Copied";
      window.setTimeout(() => {
        copyButton.textContent = "Copy regex";
      }, 1400);
    } catch (error) {
      copyButton.textContent = "Copy failed";
    }
  });

  actions.appendChild(extra);
  if (shouldCollapse) {
    const toggleButton = document.createElement("button");
    toggleButton.type = "button";
    toggleButton.className = "toggle-button";
    toggleButton.textContent = "Expand";
    toggleButton.addEventListener("click", () => {
      const isCollapsed = pattern.classList.toggle("collapsed");
      toggleButton.textContent = isCollapsed ? "Expand" : "Collapse";
    });
    actions.appendChild(toggleButton);
  }
  actions.appendChild(copyButton);

  card.appendChild(title);
  card.appendChild(badgeRow);
  card.appendChild(pattern);
  card.appendChild(example);
  card.appendChild(actions);

  return card;
};

const renderResults = (items) => {
  resultsGrid.innerHTML = "";
  const fragment = document.createDocumentFragment();
  items.forEach((item, index) => fragment.appendChild(renderCard(item, index)));
  resultsGrid.appendChild(fragment);

  resultsCount.textContent = `${items.length} pattern${items.length === 1 ? "" : "s"}`;
  emptyState.hidden = items.length !== 0;
};

const filterRecords = () => {
  const query = normalize(searchInput.value);
  const category = categorySelect.value;

  const filtered = records.filter((item) => {
    if (category !== "all" && item.group !== category) {
      return false;
    }
    if (!query) {
      return true;
    }
    return [item.name, item.group, item.regex, item.example]
      .some((value) => normalize(value).includes(query));
  });

  renderResults(filtered);
};

const debouncedFilter = debounce(filterRecords);

const startApp = (payload) => {
  const groups = payload.regular_expresions || payload.regular_expressions || [];
  records = buildRecords(groups);
  statTotal.textContent = records.length.toString();
  populateCategories(records);
  renderResults(records);
};

const loadData = async () => {
  try {
    const response = await fetch(DATA_URL);
    if (!response.ok) {
      throw new Error(`Failed to load ${DATA_URL}`);
    }
    const text = await response.text();
    const payload = window.jsyaml.load(text);
    startApp(payload || {});
  } catch (error) {
    emptyState.hidden = false;
    emptyState.querySelector("h2").textContent = "Unable to load regex.yaml.";
    emptyState.querySelector("p").textContent = "Serve the site from the repo root so the data file is reachable.";
  }
};

searchInput.addEventListener("input", debouncedFilter);
categorySelect.addEventListener("change", filterRecords);
clearButton.addEventListener("click", () => {
  searchInput.value = "";
  categorySelect.value = "all";
  filterRecords();
});

loadData();
