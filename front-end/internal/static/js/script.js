let monitoredSites = [
    { url: "https://github.com", status: "up", latency: "124ms", lastCheck: "2 mins ago" }
];

function renderSites() {
    const tableBody = document.getElementById('sites-table-body');
    if (!tableBody) return;

    tableBody.innerHTML = monitoredSites.map((site, index) => `
        <tr>
            <td><span class="status-badge ${site.status === 'up' ? 'up' : 'down'}">
                ${site.status === 'up' ? 'Online' : 'Offline'}
            </span></td>
            <td>${site.url}</td>
            <td>${site.latency}</td>
            <td>${site.lastCheck}</td>
            <td>
                <button class="btn-text" onclick="removeSite(${index})">Remove</button>
            </td>
        </tr>
    `).join('');
}

async function handleAddSite() {
    const urlInput = document.getElementById('url-input');
    const intervalInput = document.getElementById('interval-input');
    const addBtn = document.getElementById('add-btn');

    const urlValue = urlInput.value.trim();
    const intervalValue = parseInt(intervalInput.value);

    if (!urlValue || isNaN(intervalValue)) {
        alert("Fill in the URL and interval");
        return;
    }

    addBtn.disabled = true;
    addBtn.innerText = "Adding...";

    const reqData = {
        //TODO: Replace with actual user ID from session
        userID: "550e8400-e29b-41d4-a716-446655440000",
        url: urlValue,
        interval: intervalValue
    };

    try {
        const response = await fetch(`http://localhost:8080/api/v1/addSite`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(reqData)
        });

        if (!response.ok) throw new Error(`Server error: ${response.status}`);

        monitoredSites.push({
            url: urlValue,
            status: "up",
            latency: "waiting",
            lastCheck: "Just now"
        });

        renderSites();
        urlInput.value = '';

    } catch (error) {
        console.error('Error adding site:', error);
        alert("Failed to add site (check console)");
    } finally {
        addBtn.disabled = false;
        addBtn.innerText = "Add Site";
    }
}

function removeSite(index) {
    if (confirm("Remove this site from monitoring?")) {
        monitoredSites.splice(index, 1);
        renderSites();
    }
}

document.getElementById('add-btn').addEventListener('click', handleAddSite);

renderSites();