document.addEventListener("DOMContentLoaded", () => {
    checkHealth();
    setInterval(checkHealth, 10000);

    document.getElementById("addToQueueBtn").addEventListener("click", addJobToStaging);
    document.getElementById("uploadBatchBtn").addEventListener("click", uploadAllStagedJobs);
});

let stagedJobs = [];
let totalExecutedJobs = 0;

// --- 1. Health Check ---
async function checkHealth() {
    const dot = document.getElementById("healthDot");
    const text = document.getElementById("healthText");

    try {
        const response = await fetch("/checkHealth");
        if (response.ok) {
            dot.className = "w-2.5 h-2.5 rounded-full bg-emerald-400 shadow-[0_0_8px_#34d399]";
            text.textContent = "Server Online";
        } else {
            throw new Error("Unhealthy");
        }
    } catch (err) {
        dot.className = "w-2.5 h-2.5 rounded-full bg-rose-500 shadow-[0_0_8px_#f43f5e]";
        text.textContent = "Server Offline";
    }
}

// --- 2. Staging Batch Management ---
function addJobToStaging() {
    const fileInput = document.getElementById("jobFile");
    const jobTypeSelect = document.getElementById("jobType");

    if (!fileInput.files[0]) {
        alert("Please select a video file first.");
        return;
    }

    const file = fileInput.files[0];
    const jobType = jobTypeSelect.value;
    const tempId = "staged-" + Date.now() + "-" + Math.random().toString(36).substr(2, 4);

    stagedJobs.push({
        tempId: tempId,
        jobType: jobType,
        file: file
    });

    // Reset file input for next pick
    fileInput.value = "";

    renderStagedList();
}

function removeStagedJob(tempId) {
    stagedJobs = stagedJobs.filter(item => item.tempId !== tempId);
    renderStagedList();
}

function renderStagedList() {
    const stagedList = document.getElementById("stagedList");
    const stagedCount = document.getElementById("stagedCount");
    const uploadBatchBtn = document.getElementById("uploadBatchBtn");
    const uploadBtnText = document.getElementById("uploadBtnText");

    stagedCount.textContent = `${stagedJobs.length} ${stagedJobs.length === 1 ? 'Item' : 'Items'}`;
    uploadBtnText.textContent = `Upload All Jobs (${stagedJobs.length})`;

    if (stagedJobs.length === 0) {
        uploadBatchBtn.disabled = true;
        stagedList.innerHTML = `
            <div id="emptyStagedState" class="text-center py-6 text-gray-500 bg-gray-950/40 rounded-lg border border-dashed border-emerald-900/20">
                No tasks added to batch yet.
            </div>`;
        return;
    }

    uploadBatchBtn.disabled = false;
    stagedList.innerHTML = "";

    stagedJobs.forEach((item, index) => {
        const row = document.createElement("div");
        row.className = "flex justify-between items-center bg-gray-950 p-2.5 rounded-lg border border-emerald-900/40";
        row.innerHTML = `
            <div class="truncate mr-2">
                <span class="text-[10px] font-bold text-emerald-400 bg-emerald-950 px-1.5 py-0.5 rounded border border-emerald-800/60 uppercase mr-1.5">
                    ${item.jobType}
                </span>
                <span class="text-white font-medium text-xs truncate inline-block align-bottom max-w-[170px]">
                    ${item.file.name}
                </span>
            </div>
            <button type="button" onclick="removeStagedJob('${item.tempId}')" class="text-gray-500 hover:text-rose-400 p-1 transition-colors">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
            </button>
        `;
        stagedList.appendChild(row);
    });
}

// --- 3. Parallel Batch Execution ---
async function uploadAllStagedJobs() {
    if (stagedJobs.length === 0) return;

    const uploadBatchBtn = document.getElementById("uploadBatchBtn");
    const uploadBtnText = document.getElementById("uploadBtnText");

    // Copy queue reference and clear staging list
    const jobsToUpload = [...stagedJobs];
    stagedJobs = [];
    renderStagedList();

    uploadBatchBtn.disabled = true;
    uploadBtnText.textContent = "Processing Uploads...";

    // Dispatch all requests concurrently
    const uploadPromises = jobsToUpload.map(async (task) => {
        const formData = new FormData();
        formData.append("jobType", task.jobType);
        formData.append("jobFile", task.file);

        try {
            const response = await fetch("/acceptJob", {
                method: "POST",
                body: formData,
            });

            if (!response.ok) throw new Error("Job upload failed");

            const data = await response.json();
            const jobId = data.JobUniqueId || data.jobId || data;

            if (!jobId || typeof jobId !== "string") throw new Error("Invalid Job ID received");

            createJobCard(jobId, task.jobType, task.file.name);
            pollJobStatus(jobId);
        } catch (err) {
            alert(`Error uploading ${task.file.name}: ${err.message}`);
        }
    });

    await Promise.all(uploadPromises);
}

// --- 4. Executed Job Card Builder ---
function createJobCard(jobId, jobType, fileName) {
    const emptyState = document.getElementById("emptyState");
    const jobsList = document.getElementById("jobsList");
    const jobCount = document.getElementById("jobCount");

    if (emptyState) emptyState.classList.add("hidden");

    totalExecutedJobs++;
    jobCount.textContent = `${totalExecutedJobs} Executed`;

    const card = document.createElement("div");
    card.id = `job-card-${jobId}`;
    card.className = "bg-gray-900 p-5 rounded-2xl border border-emerald-900/50 shadow-xl space-y-3 transition-all";

    card.innerHTML = `
        <div class="flex justify-between items-start">
            <div>
                <span class="text-xs font-bold text-emerald-400 uppercase tracking-wider">${jobType}</span>
                <h3 class="text-sm font-semibold text-white truncate max-w-xs">${fileName}</h3>
            </div>
            <span id="badge-${jobId}" class="px-2.5 py-1 rounded-md text-xs font-bold uppercase tracking-wider bg-amber-500/10 text-amber-400 border border-amber-500/20 animate-pulse">
                PENDING
            </span>
        </div>

        <div class="text-xs text-gray-400 font-mono bg-gray-950 p-2.5 rounded-xl border border-emerald-900/30 truncate">
            ID: <span class="text-emerald-200">${jobId}</span>
        </div>

        <!-- Download Button Container -->
        <div id="download-container-${jobId}" class="hidden pt-1">
            <button onclick="triggerDownload('${jobId}')" class="w-full bg-emerald-600 hover:bg-emerald-500 text-white text-xs font-bold py-2.5 px-4 rounded-xl transition-all shadow-md flex justify-center items-center space-x-2">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path></svg>
                <span>Download Result (.zip)</span>
            </button>
        </div>
    `;

    jobsList.prepend(card);
}

// --- 5. Status Polling Loop ---
function pollJobStatus(jobId) {
    const interval = setInterval(async () => {
        try {
            const response = await fetch(`/getJobDetail/${jobId}`);
            if (!response.ok) return;

            const data = await response.json();
            const status = (data.Status || data.status || "").toUpperCase();

            updateJobBadge(jobId, status);

            if (status === "DONE" || status === "COMPLETED") {
                clearInterval(interval);
                showDownloadButton(jobId);
            } else if (status === "FAILED") {
                clearInterval(interval);
            }
        } catch (err) {
            console.error(`Error polling job ${jobId}:`, err);
        }
    }, 1000);
}

// --- 6. Helper Functions ---
function updateJobBadge(jobId, status) {
    const badge = document.getElementById(`badge-${jobId}`);
    if (!badge) return;

    if (status === "PROCESSING" || status === "PENDING") {
        badge.className = "px-2.5 py-1 rounded-md text-xs font-bold uppercase tracking-wider bg-amber-500/10 text-amber-400 border border-amber-500/20 animate-pulse";
        badge.textContent = status;
    } else if (status === "DONE" || status === "COMPLETED") {
        badge.className = "px-2.5 py-1 rounded-md text-xs font-bold uppercase tracking-wider bg-emerald-500/20 text-emerald-400 border border-emerald-500/40";
        badge.textContent = "DONE";
    } else if (status === "FAILED") {
        badge.className = "px-2.5 py-1 rounded-md text-xs font-bold uppercase tracking-wider bg-rose-500/20 text-rose-400 border border-rose-500/40";
        badge.textContent = "FAILED";
    }
}

function showDownloadButton(jobId) {
    const container = document.getElementById(`download-container-${jobId}`);
    if (container) {
        container.classList.remove("hidden");
    }
}

function triggerDownload(jobId) {
    window.location.href = `/download/${jobId}`;
}
