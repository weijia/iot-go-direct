/**
 * ConfigManager - Virtual Glass Device Configuration Manager
 *
 * Features:
 * - Load default config from config.json
 * - Override with localStorage saved config
 * - Export/import config as JSON file
 * - Validate config schema
 * - Provide typed access to config values
 */

const STORAGE_KEY = 'virtual-glass-config';

const DEFAULT_CONFIG = {
  mqtt: {
    brokerUrl: 'ws://app.kosglass.com:8083/mqtt',
    username: '',
    password: '',
    gatewayId: 'F12309150001',
    nodeId: '12345678',
    reconnectPeriod: 10000,
    connectTimeout: 30000,
    qos: 0
  },
  device: {
    hardVersion: '100',
    softVersion: '101',
    runArea: 1,
    rssi: -50,
    snr: 8,
    isOffline: false,
    completionStatus: 2
  },
  display: {
    glassSize: 320,
    autoConnect: false,
    logMaxEntries: 200
  },
  topics: {
    inbound: 'device/{gatewayId}/in',
    outbound: 'device/{gatewayId}/out'
  },
  configVersion: '1.0.0'
};

class ConfigManager {
  constructor() {
    this.config = JSON.parse(JSON.stringify(DEFAULT_CONFIG));
    this.listeners = [];
    this.loaded = false;
  }

  /**
   * Initialize config: try config.json first, then localStorage, then defaults
   */
  async init() {
    try {
      // Try to load from config.json
      const response = await fetch('config.json', { cache: 'no-store' });
      if (response.ok) {
        const fileConfig = await response.json();
        this.merge(fileConfig);
      }
    } catch (e) {
      console.warn('[Config] Failed to load config.json, using defaults:', e.message);
    }

    // Override with localStorage if present
    try {
      const stored = localStorage.getItem(STORAGE_KEY);
      if (stored) {
        const localConfig = JSON.parse(stored);
        this.merge(localConfig);
      }
    } catch (e) {
      console.warn('[Config] Failed to load from localStorage:', e.message);
    }

    this.loaded = true;
    this.notify();
    return this.config;
  }

  /**
   * Deep merge config objects
   */
  merge(source) {
    const deepMerge = (target, src) => {
      for (const key of Object.keys(src)) {
        if (src[key] && typeof src[key] === 'object' && !Array.isArray(src[key])) {
          if (!target[key] || typeof target[key] !== 'object') {
            target[key] = {};
          }
          deepMerge(target[key], src[key]);
        } else {
          target[key] = src[key];
        }
      }
    };
    deepMerge(this.config, source);
  }

  /**
   * Get config value by dot-notation path
   * e.g., get('mqtt.brokerUrl')
   */
  get(path, defaultValue = undefined) {
    const parts = path.split('.');
    let current = this.config;
    for (const part of parts) {
      if (current == null || current[part] === undefined) {
        return defaultValue;
      }
      current = current[part];
    }
    return current;
  }

  /**
   * Set config value by dot-notation path
   */
  set(path, value) {
    const parts = path.split('.');
    let current = this.config;
    for (let i = 0; i < parts.length - 1; i++) {
      if (!current[parts[i]] || typeof current[parts[i]] !== 'object') {
        current[parts[i]] = {};
      }
      current = current[parts[i]];
    }
    current[parts[parts.length - 1]] = value;
    this.notify();
  }

  /**
   * Get full config (read-only)
   */
  getAll() {
    return JSON.parse(JSON.stringify(this.config));
  }

  /**
   * Save current config to localStorage
   */
  save() {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(this.config));
      return true;
    } catch (e) {
      console.error('[Config] Failed to save to localStorage:', e.message);
      return false;
    }
  }

  /**
   * Reset to defaults (clear localStorage and reload from config.json)
   */
  async reset() {
    localStorage.removeItem(STORAGE_KEY);
    this.config = JSON.parse(JSON.stringify(DEFAULT_CONFIG));
    await this.init();
  }

  /**
   * Export config as JSON file download
   */
  export(filename = 'virtual-glass-config.json') {
    const dataStr = JSON.stringify(this.config, null, 2);
    const blob = new Blob([dataStr], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  }

  /**
   * Import config from a JSON file
   */
  import(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        try {
          const imported = JSON.parse(e.target.result);
          this.config = JSON.parse(JSON.stringify(DEFAULT_CONFIG));
          this.merge(imported);
          this.save();
          this.notify();
          resolve(this.config);
        } catch (err) {
          reject(new Error('Invalid JSON file: ' + err.message));
        }
      };
      reader.onerror = () => reject(new Error('Failed to read file'));
      reader.readAsText(file);
    });
  }

  /**
   * Get MQTT topic with placeholders replaced
   */
  getTopic(type) {
    const template = this.config.topics[type] || '';
    return template.replace('{gatewayId}', this.config.mqtt.gatewayId);
  }

  /**
   * Subscribe to config changes
   */
  subscribe(callback) {
    this.listeners.push(callback);
    return () => {
      this.listeners = this.listeners.filter(l => l !== callback);
    };
  }

  /**
   * Notify all listeners of config change
   */
  notify() {
    for (const listener of this.listeners) {
      try {
        listener(this.config);
      } catch (e) {
        console.error('[Config] Listener error:', e);
      }
    }
  }

  /**
   * Validate current config
   */
  validate() {
    const errors = [];
    const mqtt = this.config.mqtt;

    if (!mqtt.brokerUrl || !mqtt.brokerUrl.startsWith('ws')) {
      errors.push('Broker URL must start with ws:// or wss://');
    }
    if (!mqtt.gatewayId || mqtt.gatewayId.length < 8) {
      errors.push('Gateway Node ID is required (8+ hex chars)');
    }
    if (!mqtt.nodeId || mqtt.nodeId.length < 4) {
      errors.push('Virtual Node ID is required');
    }

    return {
      valid: errors.length === 0,
      errors
    };
  }
}

// Singleton instance
const configManager = new ConfigManager();
