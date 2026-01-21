let monitoredSites = [];
let sitesRefreshInterval = null;
let userID = "550e8400-e29b-41d4-a716-446655440000";
function renderSites() {
    const tableBody = document.getElementById('sites-table-body');
    if (!tableBody) return;

    tableBody.innerHTML = monitoredSites.map((site, index) => `
        <tr>
            <td><span class="status-badge ${site.status}">
                ${site.status}
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
async function fetchMonitoredSites() {
    try {
        const userID = "550e8400-e29b-41d4-a716-446655440000";
        const response = await fetch(`http://localhost:8080/api/v1/getUserSites/${userID}`, {
            method: 'GET',
            headers: {'Content-Type': 'application/json'},
        });

        if (!response.ok) throw new Error(`Server error: ${response.status}`);
        const data = await response.json();

        const sites = data.sites;
        monitoredSites = []
        const tableBody = document.getElementById('sites-table-body');
        tableBody.innerHTML = "";
        for (const site of sites) {
            siteStatus = site.is_up ? "up" : "down";
            const lastCheck = new Date(site.last_check).toLocaleString();
            monitoredSites.push({
                url: site.url,
                status: siteStatus,
                latency: "",
                lastCheck: lastCheck
            });
        }

        renderSites();

    } catch (error) {
        console.error('Error fetching sites:', error);
        alert("Failed to fetch sites (check console)");
    }
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
            status: "WAITING",
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

async function removeSite(index) {
    const siteToRemove = monitoredSites[index];
    if (!siteToRemove) return;

    if (!confirm(`Remove ${siteToRemove.url} from monitoring?`)) {
        return;
    }

    const reqData = {
        userID: "550e8400-e29b-41d4-a716-446655440000",
        url: siteToRemove.url,
    };

    try {
        const response = await fetch(`http://localhost:8080/api/v1/deleteSite`, {
            method: 'DELETE',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(reqData)
        });

        if (!response.ok) {
            throw new Error(`Server error: ${response.status}`);
        }

        monitoredSites.splice(index, 1);

        renderSites();

    } catch (error) {
        console.error('Error deleting site:', error);
        alert("Failed to delete site. It might still be in the list.");
    }
}

function startMonitoring(userID, intervalMs = 5000) {
    if (sitesRefreshInterval) {
        clearInterval(sitesRefreshInterval);
    }

    fetchMonitoredSites();

    sitesRefreshInterval = setInterval(() => {
        fetchMonitoredSites();
    }, intervalMs);

}

document.getElementById('add-btn').addEventListener('click', handleAddSite);

startMonitoring(userID, 5000);

renderSites();