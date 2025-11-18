// Default configuration
const defaultConfig = {
    server: {
        port: "3030",
        host: "0.0.0.0",
        environment: "development",
        read_timeout: 10,
        write_timeout: 10
    },
    database: {
        enabled: false,
        type: "postgres",
        host: "localhost",
        port: 5432,
        database: "countries",
        username: "postgres",
        password: "",
        ssl_mode: "disable",
        schema: {
            countries_table: "countries",
            aliases_table: "country_aliases",
            code_column: "code",
            name_column: "name",
            alias_code_column: "country_code",
            alias_name_column: "alias"
        }
    },
    data: {
        source: "csv",
        countries_file: "data/countries.csv",
        aliases_file: "data/aliases.csv"
    },
    logging: {
        level: "info",
        format: "json"
    },
    gui: {
        enabled: true,
        path: "/admin"
    }
};

// Initialize page on load
document.addEventListener('DOMContentLoaded', function() {
    resetConfig();
    updateYAMLOutput();
});

// Toggle data source fields
function toggleDataSourceFields() {
    const source = document.getElementById('data-source').value;
    const fileFields = document.getElementById('file-fields');
    const dbSection = document.getElementById('database-section');

    if (source === 'csv' || source === 'tsv') {
        fileFields.style.display = 'block';
    } else {
        fileFields.style.display = 'none';
    }

    if (source === 'database') {
        dbSection.style.display = 'block';
        document.getElementById('db-enabled').checked = true;
    } else {
        dbSection.style.display = 'none';
        document.getElementById('db-enabled').checked = false;
    }

    updateYAMLOutput();
}

// Get current configuration from form
function getCurrentConfig() {
    return {
        server: {
            port: document.getElementById('server-port').value,
            host: document.getElementById('server-host').value,
            environment: document.getElementById('server-environment').value,
            read_timeout: parseInt(document.getElementById('server-read-timeout').value),
            write_timeout: parseInt(document.getElementById('server-write-timeout').value)
        },
        database: {
            enabled: document.getElementById('db-enabled').checked,
            type: document.getElementById('db-type').value,
            host: document.getElementById('db-host').value,
            port: parseInt(document.getElementById('db-port').value),
            database: document.getElementById('db-name').value,
            username: document.getElementById('db-username').value,
            password: document.getElementById('db-password').value,
            ssl_mode: document.getElementById('db-ssl-mode').value,
            schema: defaultConfig.database.schema
        },
        data: {
            source: document.getElementById('data-source').value,
            countries_file: document.getElementById('countries-file').value,
            aliases_file: document.getElementById('aliases-file').value
        },
        logging: {
            level: document.getElementById('log-level').value,
            format: document.getElementById('log-format').value
        },
        gui: {
            enabled: document.getElementById('gui-enabled').checked,
            path: document.getElementById('gui-path').value
        }
    };
}

// Load configuration into form
function loadConfigIntoForm(config) {
    document.getElementById('server-port').value = config.server.port;
    document.getElementById('server-host').value = config.server.host;
    document.getElementById('server-environment').value = config.server.environment;
    document.getElementById('server-read-timeout').value = config.server.read_timeout;
    document.getElementById('server-write-timeout').value = config.server.write_timeout;

    document.getElementById('db-enabled').checked = config.database.enabled;
    document.getElementById('db-type').value = config.database.type;
    document.getElementById('db-host').value = config.database.host;
    document.getElementById('db-port').value = config.database.port;
    document.getElementById('db-name').value = config.database.database;
    document.getElementById('db-username').value = config.database.username;
    document.getElementById('db-password').value = config.database.password;
    document.getElementById('db-ssl-mode').value = config.database.ssl_mode;

    document.getElementById('data-source').value = config.data.source;
    document.getElementById('countries-file').value = config.data.countries_file;
    document.getElementById('aliases-file').value = config.data.aliases_file;

    document.getElementById('log-level').value = config.logging.level;
    document.getElementById('log-format').value = config.logging.format;

    document.getElementById('gui-enabled').checked = config.gui.enabled;
    document.getElementById('gui-path').value = config.gui.path;

    toggleDataSourceFields();
}

// Convert config object to YAML string
function configToYAML(config) {
    return `server:
  port: "${config.server.port}"
  host: "${config.server.host}"
  environment: "${config.server.environment}"
  read_timeout: ${config.server.read_timeout}
  write_timeout: ${config.server.write_timeout}

database:
  enabled: ${config.database.enabled}
  type: "${config.database.type}"
  host: "${config.database.host}"
  port: ${config.database.port}
  database: "${config.database.database}"
  username: "${config.database.username}"
  password: "${config.database.password}"
  ssl_mode: "${config.database.ssl_mode}"
  schema:
    countries_table: "${config.database.schema.countries_table}"
    aliases_table: "${config.database.schema.aliases_table}"
    code_column: "${config.database.schema.code_column}"
    name_column: "${config.database.schema.name_column}"
    alias_code_column: "${config.database.schema.alias_code_column}"
    alias_name_column: "${config.database.schema.alias_name_column}"

data:
  source: "${config.data.source}"
  countries_file: "${config.data.countries_file}"
  aliases_file: "${config.data.aliases_file}"

logging:
  level: "${config.logging.level}"
  format: "${config.logging.format}"

gui:
  enabled: ${config.gui.enabled}
  path: "${config.gui.path}"
`;
}

// Update YAML output
function updateYAMLOutput() {
    const config = getCurrentConfig();
    const yaml = configToYAML(config);
    document.getElementById('yaml-output').textContent = yaml;
}

// Add event listeners to all inputs
document.querySelectorAll('input, select').forEach(element => {
    element.addEventListener('input', updateYAMLOutput);
    element.addEventListener('change', updateYAMLOutput);
});

// Test country lookup
async function testLookup() {
    const countryInput = document.getElementById('country-input').value.trim();
    const resultDiv = document.getElementById('lookup-result');

    if (!countryInput) {
        resultDiv.innerHTML = '<div class="error">Please enter a country name</div>';
        return;
    }

    resultDiv.innerHTML = '<div class="loading">Looking up...</div>';

    try {
        const response = await fetch(`/api/gui/lookup?country=${encodeURIComponent(countryInput)}`);
        const data = await response.json();

        if (response.ok) {
            resultDiv.innerHTML = `
                <div class="success">
                    <h3>✅ Match Found!</h3>
                    <table>
                        <tr>
                            <td><strong>Country Code:</strong></td>
                            <td>${data.isoCode}</td>
                        </tr>
                        <tr>
                            <td><strong>Official Name:</strong></td>
                            <td>${data.officialName}</td>
                        </tr>
                        <tr>
                            <td><strong>Query:</strong></td>
                            <td>${data.query}</td>
                        </tr>
                    </table>
                </div>
            `;
        } else {
            resultDiv.innerHTML = `
                <div class="error">
                    <h3>❌ ${data.error || 'Error'}</h3>
                    <p>${data.message || 'Country not found'}</p>
                </div>
            `;
        }
    } catch (error) {
        resultDiv.innerHTML = `
            <div class="error">
                <h3>❌ Request Failed</h3>
                <p>${error.message}</p>
            </div>
        `;
    }
}

// Load statistics
async function loadStats() {
    const statsDiv = document.getElementById('stats-result');
    statsDiv.innerHTML = '<div class="loading">Loading statistics...</div>';

    try {
        const response = await fetch('/api/gui/stats');
        const data = await response.json();

        if (response.ok) {
            const successRate = (data.success_rate * 100).toFixed(2);
            const failureRate = (data.failure_rate * 100).toFixed(2);

            let popularCountriesHTML = '';
            if (data.popular_countries && data.popular_countries.length > 0) {
                popularCountriesHTML = `
                    <h3>🏆 Most Popular Countries</h3>
                    <table>
                        <tr>
                            <th>Rank</th>
                            <th>Code</th>
                            <th>Name</th>
                            <th>Count</th>
                        </tr>
                        ${data.popular_countries.map((country, index) => `
                            <tr>
                                <td>${index + 1}</td>
                                <td>${country.code}</td>
                                <td>${country.name}</td>
                                <td>${country.count}</td>
                            </tr>
                        `).join('')}
                    </table>
                `;
            }

            statsDiv.innerHTML = `
                <div class="stats-content">
                    <div class="stats-grid">
                        <div class="stat-card">
                            <div class="stat-value">${data.total_requests}</div>
                            <div class="stat-label">Total Requests</div>
                        </div>
                        <div class="stat-card success">
                            <div class="stat-value">${data.success_count}</div>
                            <div class="stat-label">Successful</div>
                        </div>
                        <div class="stat-card error">
                            <div class="stat-value">${data.not_found_count}</div>
                            <div class="stat-label">Not Found</div>
                        </div>
                        <div class="stat-card error">
                            <div class="stat-value">${data.error_count}</div>
                            <div class="stat-label">Errors</div>
                        </div>
                    </div>
                    <div class="stats-rates">
                        <div class="rate-bar">
                            <div class="rate-label">Success Rate: ${successRate}%</div>
                            <div class="rate-progress">
                                <div class="rate-fill success" style="width: ${successRate}%"></div>
                            </div>
                        </div>
                        <div class="rate-bar">
                            <div class="rate-label">Failure Rate: ${failureRate}%</div>
                            <div class="rate-progress">
                                <div class="rate-fill error" style="width: ${failureRate}%"></div>
                            </div>
                        </div>
                    </div>
                    ${popularCountriesHTML}
                </div>
            `;
        } else {
            statsDiv.innerHTML = `
                <div class="error">
                    <h3>❌ Failed to Load Statistics</h3>
                    <p>${data.message || 'Unknown error'}</p>
                </div>
            `;
        }
    } catch (error) {
        statsDiv.innerHTML = `
            <div class="error">
                <h3>❌ Request Failed</h3>
                <p>${error.message}</p>
            </div>
        `;
    }
}

// Save configuration
async function saveConfig() {
    const config = getCurrentConfig();

    try {
        const response = await fetch('/api/config/save', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(config)
        });

        const data = await response.json();

        if (response.ok) {
            showMessage('Configuration saved successfully!', 'success');
        } else {
            showMessage('Failed to save: ' + (data.message || 'Unknown error'), 'error');
        }
    } catch (error) {
        showMessage('Failed to save configuration: ' + error.message, 'error');
    }
}

// Load configuration from server
async function loadConfig() {
    try {
        const response = await fetch('/api/config');

        if (response.ok) {
            const config = await response.json();
            loadConfigIntoForm(config);
            updateYAMLOutput();
            showMessage('Configuration loaded from server', 'success');
        } else {
            showMessage('Failed to load configuration from server', 'error');
            loadConfigIntoForm(defaultConfig);
            updateYAMLOutput();
        }
    } catch (error) {
        showMessage('Failed to load configuration: ' + error.message, 'error');
        loadConfigIntoForm(defaultConfig);
        updateYAMLOutput();
    }
}

// Download configuration as YAML file
function downloadConfig() {
    const config = getCurrentConfig();
    const yaml = configToYAML(config);
    const blob = new Blob([yaml], { type: 'text/yaml' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'config.yaml';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    showMessage('Configuration downloaded successfully!', 'success');
}

// Reset to default configuration
function resetConfig() {
    loadConfigIntoForm(defaultConfig);
    updateYAMLOutput();
    showMessage('Configuration reset to defaults.', 'success');
}

// Show message
function showMessage(text, type) {
    const messageEl = document.getElementById('message');
    messageEl.textContent = text;
    messageEl.className = 'message ' + type;

    setTimeout(() => {
        messageEl.className = 'message';
    }, 5000);
}

// Bulk processing variables
let uploadedFile = null;
let bulkResults = [];

// Handle file selection
function handleFileSelect() {
    const fileInput = document.getElementById('file-input');
    const processBtn = document.getElementById('process-btn');

    if (fileInput.files.length > 0) {
        uploadedFile = fileInput.files[0];
        processBtn.disabled = false;

        const resultDiv = document.getElementById('bulk-result');
        resultDiv.innerHTML = `<div class="info">File selected: ${uploadedFile.name} (${(uploadedFile.size / 1024).toFixed(2)} KB)</div>`;
    } else {
        uploadedFile = null;
        processBtn.disabled = true;
    }
}

// Parse CSV/TSV content
function parseFile(content, delimiter) {
    const lines = content.trim().split('\n');
    if (lines.length === 0) {
        return { headers: [], rows: [] };
    }

    const headers = lines[0].split(delimiter).map(h => h.trim().replace(/^["']|["']$/g, ''));
    const rows = [];

    for (let i = 1; i < lines.length; i++) {
        if (lines[i].trim() === '') continue;

        const values = lines[i].split(delimiter).map(v => v.trim().replace(/^["']|["']$/g, ''));
        const row = {};
        headers.forEach((header, index) => {
            row[header] = values[index] || '';
        });
        rows.push(row);
    }

    return { headers, rows };
}

// Process bulk lookup
async function processBulk() {
    if (!uploadedFile) {
        return;
    }

    const resultDiv = document.getElementById('bulk-result');
    const processBtn = document.getElementById('process-btn');
    const downloadBtn = document.getElementById('download-btn');
    const columnName = document.getElementById('column-name').value.trim();

    processBtn.disabled = true;
    downloadBtn.style.display = 'none';
    resultDiv.innerHTML = '<div class="loading">Reading file...</div>';

    try {
        const content = await uploadedFile.text();
        const delimiter = uploadedFile.name.endsWith('.tsv') ? '\t' : ',';
        const { headers, rows } = parseFile(content, delimiter);

        if (rows.length === 0) {
            resultDiv.innerHTML = '<div class="error">No data found in file</div>';
            processBtn.disabled = false;
            return;
        }

        // Find the column index
        let countryColumn = columnName;
        if (!headers.includes(countryColumn)) {
            // Try common alternatives
            const alternatives = ['country', 'country_name', 'name', 'Country', 'Country Name'];
            countryColumn = alternatives.find(alt => headers.includes(alt));

            if (!countryColumn) {
                resultDiv.innerHTML = `<div class="error">Column "${columnName}" not found. Available columns: ${headers.join(', ')}</div>`;
                processBtn.disabled = false;
                return;
            }
        }

        resultDiv.innerHTML = `<div class="loading">Processing ${rows.length} countries...</div>`;

        // Process lookups
        bulkResults = [];
        let processed = 0;

        for (const row of rows) {
            const countryName = row[countryColumn];
            if (!countryName || countryName.trim() === '') {
                bulkResults.push({
                    query: '',
                    isoCode: '',
                    officialName: '',
                    status: 'skipped',
                    error: 'Empty country name'
                });
                continue;
            }

            try {
                const response = await fetch(`/api/gui/lookup?country=${encodeURIComponent(countryName)}`);
                const data = await response.json();

                if (response.ok) {
                    bulkResults.push({
                        query: data.query,
                        isoCode: data.isoCode,
                        officialName: data.officialName,
                        status: 'success'
                    });
                } else {
                    bulkResults.push({
                        query: countryName,
                        isoCode: '',
                        officialName: '',
                        status: 'not_found',
                        error: data.message || 'Not found'
                    });
                }
            } catch (error) {
                bulkResults.push({
                    query: countryName,
                    isoCode: '',
                    officialName: '',
                    status: 'error',
                    error: error.message
                });
            }

            processed++;
            if (processed % 10 === 0) {
                resultDiv.innerHTML = `<div class="loading">Processing ${processed}/${rows.length} countries...</div>`;
            }
        }

        displayBulkResults();
        processBtn.disabled = false;
        downloadBtn.style.display = 'inline-block';
        downloadBtn.disabled = false;

    } catch (error) {
        resultDiv.innerHTML = `<div class="error">Failed to process file: ${error.message}</div>`;
        processBtn.disabled = false;
    }
}

// Display bulk results
function displayBulkResults() {
    const resultDiv = document.getElementById('bulk-result');

    const successCount = bulkResults.filter(r => r.status === 'success').length;
    const notFoundCount = bulkResults.filter(r => r.status === 'not_found').length;
    const errorCount = bulkResults.filter(r => r.status === 'error').length;
    const skippedCount = bulkResults.filter(r => r.status === 'skipped').length;

    let tableHTML = `
        <div class="success">
            <h3>✅ Bulk Processing Complete</h3>
            <p>
                Total: ${bulkResults.length} |
                Success: ${successCount} |
                Not Found: ${notFoundCount} |
                Error: ${errorCount} |
                Skipped: ${skippedCount}
            </p>
            <table>
                <thead>
                    <tr>
                        <th>#</th>
                        <th>Query</th>
                        <th>ISO Code</th>
                        <th>Official Name</th>
                        <th>Status</th>
                    </tr>
                </thead>
                <tbody>
    `;

    bulkResults.forEach((result, index) => {
        const statusClass = result.status === 'success' ? 'success' : 'error';
        const statusIcon = result.status === 'success' ? '✅' : '❌';
        const displayStatus = result.status === 'success' ? 'Success' : (result.error || result.status);

        tableHTML += `
            <tr class="${statusClass}">
                <td>${index + 1}</td>
                <td>${result.query}</td>
                <td>${result.isoCode || '-'}</td>
                <td>${result.officialName || '-'}</td>
                <td>${statusIcon} ${displayStatus}</td>
            </tr>
        `;
    });

    tableHTML += `
                </tbody>
            </table>
        </div>
    `;

    resultDiv.innerHTML = tableHTML;
}

// Download results as CSV
function downloadResults() {
    if (bulkResults.length === 0) {
        return;
    }

    let csv = 'Query,ISO Code,Official Name,Status,Error\n';

    bulkResults.forEach(result => {
        const row = [
            `"${result.query}"`,
            `"${result.isoCode || ''}"`,
            `"${result.officialName || ''}"`,
            `"${result.status}"`,
            `"${result.error || ''}"`
        ];
        csv += row.join(',') + '\n';
    });

    const blob = new Blob([csv], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'bulk-lookup-results.csv';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
}
