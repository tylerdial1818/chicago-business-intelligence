'use client';

import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import type { ZipCode } from '@/types';

interface ZipCodeSelectorProps {
  value: string;
  onChange: (zipCode: string) => void;
}

export function ZipCodeSelector({ value, onChange }: ZipCodeSelectorProps) {
  const [zipCodes, setZipCodes] = useState<ZipCode[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.zipCodes()
      .then(setZipCodes)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="animate-pulse">
        <div className="h-10 bg-slate-700 rounded"></div>
      </div>
    );
  }

  return (
    <select
      value={value}
      onChange={(e) => onChange(e.target.value)}
      className="w-full px-4 py-2 bg-slate-700 border border-slate-600 rounded-lg text-slate-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
    >
      <option value="">Select ZIP Code</option>
      {zipCodes.map((zip) => (
        <option key={zip.zip_code} value={zip.zip_code}>
          {zip.zip_code} - {zip.neighborhood}
        </option>
      ))}
    </select>
  );
}
