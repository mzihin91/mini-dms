import axios from 'axios';

// Create axios instance with base URL
// In Docker: nginx proxies /api/* to backend, so use relative path
// For local dev: use environment variable or localhost
const apiClient = axios.create({
  baseURL: process.env.REACT_APP_API_URL || '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for logging
apiClient.interceptors.request.use(
  (config) => {
    console.log(`API Request: ${config.method?.toUpperCase()} ${config.url}`);
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    console.error('API Error:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);

export default apiClient;

// API functions

// Health check
export const checkHealth = async () => {
  const response = await apiClient.get('/health');
  return response.data;
};

// Device APIs
export const fetchDevices = async () => {
  const response = await apiClient.get('/devices');
  return response.data;
};

export const fetchDeviceById = async (id) => {
  const response = await apiClient.get(`/devices/${id}`);
  return response.data;
};

export const createDevice = async (deviceData) => {
  const response = await apiClient.post('/devices', deviceData);
  return response.data;
};

export const deleteDevice = async (id) => {
  const response = await apiClient.delete(`/devices/${id}`);
  return response.data;
};

export const activateDevice = async (id) => {
  const response = await apiClient.post(`/devices/${id}/activate`);
  return response.data;
};

export const deactivateDevice = async (id) => {
  const response = await apiClient.post(`/devices/${id}/deactivate`);
  return response.data;
};

// Transaction APIs
export const fetchTransactions = async (params = {}) => {
  const response = await apiClient.get('/transactions', { params });
  return response.data;
};
