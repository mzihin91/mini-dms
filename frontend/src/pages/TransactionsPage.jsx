import React, { useState, useEffect } from 'react';
import TransactionList from '../components/TransactionList';
import { fetchTransactions as fetchTransactionsAPI } from '../services/api';

function TransactionsPage() {
  const [transactions, setTransactions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const loadTransactions = async () => {
    try {
      setLoading(true);
      const data = await fetchTransactionsAPI();
      setTransactions(data.transactions || []);
      setError('');
    } catch (err) {
      setError(err.response?.data?.error || err.message);
      console.error('Error fetching transactions:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    // Initial fetch
    loadTransactions();
  }, []);

  return (
    <div className="px-4 py-6 sm:px-0">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Transaction Monitor</h1>
        <button
          onClick={loadTransactions}
          disabled={loading}
          className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors duration-200"
        >
          {loading ? 'Refreshing...' : 'Refresh'}
        </button>
      </div>
      
      {error && (
        <div className="mb-4 p-4 bg-red-50 border border-red-200 text-red-700 rounded-lg">
          {error}
        </div>
      )}

      <TransactionList transactions={transactions} loading={loading} />
    </div>
  );
}

export default TransactionsPage;
