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

// Save configuration
function saveConfig() {
    const config = getCurrentConfig();
    const yaml = configToYAML(config);

    // In a real implementation, this would send to the server
    // For now, we'll just show a success message and offer download
    showMessage('Configuration ready! Click "Download YAML" to save to file.', 'success');
    console.log('Configuration to save:', config);
}

// Load configuration from server
async function loadConfig() {
    try {
        // In a real implementation, this would fetch from the server
        // For now, we'll use the default config
        showMessage('Loading default configuration...', 'success');
        loadConfigIntoForm(defaultConfig);
        updateYAMLOutput();
    } catch (error) {
        showMessage('Failed to load configuration: ' + error.message, 'error');
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
