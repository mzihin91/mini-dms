import React, { useState } from 'react';

function TransactionList({ transactions, loading }) {
  const [expandedRows, setExpandedRows] = useState(new Set());

  const toggleRow = (transactionId) => {
    const newExpanded = new Set(expandedRows);
    if (newExpanded.has(transactionId)) {
      newExpanded.delete(transactionId);
    } else {
      newExpanded.add(transactionId);
    }
    setExpandedRows(newExpanded);
  };

  if (loading) {
    return (
      <div className="bg-white shadow-sm rounded-lg p-6">
        <div className="flex justify-center items-center py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      </div>
    );
  }

  if (!transactions || transactions.length === 0) {
    return (
      <div className="bg-white shadow-sm rounded-lg p-6">
        <p className="text-gray-500 text-center py-8">No transactions generated yet. Activate devices to see transactions.</p>
      </div>
    );
  }

  const getEventTypeBadge = (eventType) => {
    const badges = {
      access_granted: 'bg-green-100 text-green-800',
      access_denied: 'bg-red-100 text-red-800',
      face_match: 'bg-blue-100 text-blue-800',
      plate_read: 'bg-purple-100 text-purple-800',
    };

    const labels = {
      access_granted: 'Access Granted',
      access_denied: 'Access Denied',
      face_match: 'Face Match',
      plate_read: 'Plate Read',
    };

    return (
      <span className={`px-2 py-1 text-xs font-medium rounded-full ${badges[eventType] || 'bg-gray-100 text-gray-800'}`}>
        {labels[eventType] || eventType}
      </span>
    );
  };

  // Sort transactions by timestamp (most recent first)
  const sortedTransactions = [...transactions].sort((a, b) => {
    return new Date(b.timestamp) - new Date(a.timestamp);
  });

  return (
    <div className="bg-white shadow-sm rounded-lg overflow-hidden">
      <div className="px-6 py-4 border-b border-gray-200">
        <h2 className="text-lg font-semibold text-gray-900">Recent Transactions</h2>
        <p className="text-sm text-gray-500 mt-1">Showing transactions (most recent first)</p>
      </div>
      
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Device ID
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Timestamp
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Username
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Event Type
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Payload
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {sortedTransactions.map((transaction) => (
              <React.Fragment key={transaction.id}>
                <tr className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm font-medium text-gray-900">{transaction.device_id}</div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm text-gray-900">
                      {new Date(transaction.timestamp).toLocaleString()}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm text-gray-900">
                      {transaction.username || '-'}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    {getEventTypeBadge(transaction.event_type)}
                  </td>
                  <td className="px-6 py-4">
                    {transaction.payload ? (
                      <button
                        onClick={() => toggleRow(transaction.id)}
                        className="text-blue-600 hover:text-blue-800 text-sm font-medium focus:outline-none"
                      >
                        {expandedRows.has(transaction.id) ? '▼ Hide Details' : '▶ Show Details'}
                      </button>
                    ) : (
                      <span className="text-sm text-gray-500">-</span>
                    )}
                  </td>
                </tr>
                {expandedRows.has(transaction.id) && transaction.payload && (
                  <tr>
                    <td colSpan="5" className="px-6 py-4 bg-gray-50">
                      <div className="text-sm">
                        <div className="font-medium text-gray-700 mb-2">Payload Details:</div>
                        <pre className="bg-white border border-gray-200 rounded p-3 text-xs overflow-x-auto">
                          {JSON.stringify(transaction.payload, null, 2)}
                        </pre>
                      </div>
                    </td>
                  </tr>
                )}
              </React.Fragment>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default TransactionList;

