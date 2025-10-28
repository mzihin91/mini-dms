import React, { useState } from 'react';
import { createDevice } from '../services/api';

function DeviceForm({ onDeviceCreated }) {
  const [formData, setFormData] = useState({
    name: '',
    device_type: 'access_controller',
    ip_address: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const device = await createDevice(formData);
      
      // Reset form
      setFormData({
        name: '',
        device_type: 'access_controller',
        ip_address: '',
      });

      // Notify parent component
      if (onDeviceCreated) {
        onDeviceCreated(device);
      }
    } catch (err) {
      setError(err.response?.data?.error || err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  return (
    <div className="bg-white shadow-sm rounded-lg p-6 mb-6">
      <h2 className="text-lg font-semibold text-gray-900 mb-4">Register New Device</h2>
      
      {error && (
        <div className="mb-4 p-3 bg-red-50 border border-red-200 text-red-700 rounded">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-1">
            Device Name
          </label>
          <input
            type="text"
            id="name"
            name="name"
            value={formData.name}
            onChange={handleChange}
            required
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="e.g., Front Door AC"
          />
        </div>

        <div>
          <label htmlFor="device_type" className="block text-sm font-medium text-gray-700 mb-1">
            Device Type
          </label>
          <select
            id="device_type"
            name="device_type"
            value={formData.device_type}
            onChange={handleChange}
            required
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="access_controller">Access Controller</option>
            <option value="face_reader">Face Recognition Reader</option>
            <option value="anpr">ANPR Camera</option>
          </select>
        </div>

        <div>
          <label htmlFor="ip_address" className="block text-sm font-medium text-gray-700 mb-1">
            IP Address
          </label>
          <input
            type="text"
            id="ip_address"
            name="ip_address"
            value={formData.ip_address}
            onChange={handleChange}
            required
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="e.g., 192.168.1.10"
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
        >
          {loading ? 'Creating...' : 'Create Device'}
        </button>
      </form>
    </div>
  );
}

export default DeviceForm;
