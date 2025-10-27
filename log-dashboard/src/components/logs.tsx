import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { useState } from 'react';

interface Log {
  timestamp: string;
  content: string;
  is_verified: boolean;
  client_name: string;
}

const sampleLogs: Log[] = [
  { timestamp: '2023-10-01T10:00:00Z', content: 'User logged in', is_verified: true, client_name: 'Client A' },
  { timestamp: '2023-10-02T11:00:00Z', content: 'Error occurred', is_verified: false, client_name: 'Client B' },
  { timestamp: '2023-10-03T12:00:00Z', content: 'Data updated', is_verified: true, client_name: 'Client A' },
];

function Logs() {
  const [contentFilter, setContentFilter] = useState('');
  const [isVerifiedFilter, setIsVerifiedFilter] = useState<boolean | null>(null);
  const [clientNameFilter, setClientNameFilter] = useState('');
  const [dateFilter, setDateFilter] = useState('');

  const filteredLogs = sampleLogs.filter(log => {
    const matchesContent = log.content.toLowerCase().includes(contentFilter.toLowerCase());
    const matchesVerified = isVerifiedFilter === null || log.is_verified === isVerifiedFilter;
    const matchesClient = log.client_name.toLowerCase().includes(clientNameFilter.toLowerCase());
    const matchesDate = !dateFilter || log.timestamp.startsWith(dateFilter);
    return matchesContent && matchesVerified && matchesClient && matchesDate;
  });

  const clearFilters = () => {
    setContentFilter('');
    setIsVerifiedFilter(null);
    setClientNameFilter('');
    setDateFilter('');
  };

  return (
    <div className="max-w-7xl">
      <div className="flex flex-wrap gap-4 mb-4">
        <Input
          type="text"
          placeholder="Search log content"
          value={contentFilter}
          onChange={(e) => setContentFilter(e.target.value)}
          className="flex-1 min-w-0"
        />
        <Select value={isVerifiedFilter === null ? '' : isVerifiedFilter.toString()} onValueChange={(value) => setIsVerifiedFilter(value === '' ? null : value === 'true')}>
          <SelectTrigger className="w-40">
            <SelectValue placeholder="All" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="true">Verified</SelectItem>
            <SelectItem value="false">Not Verified</SelectItem>
          </SelectContent>
        </Select>
        <Input
          type="text"
          placeholder="Search client name"
          value={clientNameFilter}
          onChange={(e) => setClientNameFilter(e.target.value)}
          className="flex-1 min-w-0"
        />
        <Input
          type="date"
          value={dateFilter}
          onChange={(e) => setDateFilter(e.target.value)}
          className="w-40"
        />
        <Button onClick={clearFilters} variant="outline">Clear Filters</Button>
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Timestamp</TableHead>
            <TableHead>Content</TableHead>
            <TableHead>Verified</TableHead>
            <TableHead>Client</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {filteredLogs.map((log, index) => (
            <TableRow key={index}>
              <TableCell>{log.timestamp}</TableCell>
              <TableCell>{log.content}</TableCell>
              <TableCell>{log.is_verified ? 'Yes' : 'No'}</TableCell>
              <TableCell>{log.client_name}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}

export default Logs;