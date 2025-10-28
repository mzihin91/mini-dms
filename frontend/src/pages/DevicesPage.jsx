import React, { useState, useEffect } from 'react';
import DeviceForm from '../components/DeviceForm';
import DeviceList from '../components/DeviceList';
import { fetchDevices as fetchDevicesAPI, activateDevice, deactivateDevice } from '../services/api';

function DevicesPage() {
  const [devices, setDevices] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const loadDevices = async () => {
    try {
      setLoading(true);
      const data = await fetchDevicesAPI();
      setDevices(data.devices || []);
      setError('');
    } catch (err) {
      setError(err.response?.data?.error || err.message);
      console.error('Error fetching devices:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadDevices();
  }, []);

  const handleDeviceCreated = (newDevice) => {
    setDevices((prevDevices) => [newDevice, ...prevDevices]);
  };

  const handleActivate = async (deviceId) => {
    try {
      await activateDevice(deviceId);
      // Refresh device list
      await loadDevices();
    } catch (err) {
      setError(err.response?.data?.error || err.message);
      console.error('Error activating device:', err);
    }
  };

  const handleDeactivate = async (deviceId) => {
    try {
      await deactivateDevice(deviceId);
      // Refresh device list
      await loadDevices();
    } catch (err) {
      setError(err.response?.data?.error || err.message);
      console.error('Error deactivating device:', err);
    }
  };

  return (
    <div className="px-4 py-6 sm:px-0">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Device Management</h1>
      
      {error && (
        <div className="mb-4 p-4 bg-red-50 border border-red-200 text-red-700 rounded-lg">
          {error}
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-1">
          <DeviceForm onDeviceCreated={handleDeviceCreated} />
        </div>
        
        <div className="lg:col-span-2">
          <DeviceList 
            devices={devices} 
            loading={loading}
            onRefresh={loadDevices}
            onActivate={handleActivate}
            onDeactivate={handleDeactivate}
          />
        </div>
      </div>
    </div>
  );
}

export default DevicesPage;
