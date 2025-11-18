import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { API_BASE_URL } from "@/constants";
import { useInfiniteQuery } from "@tanstack/react-query";
import { useState } from "react";
import LogFilters from "./log-filters";
import { Button } from "./ui/button";

interface Log {
  Timestamp: string;
  Content: string;
  IsValid: boolean;
  Source: string;
}

interface LogsResponse {
  logs: Log[];
  bookmark: string | null;
  hasNextPage: boolean;
}


function Logs() {
  const [appliedFilters, setAppliedFilters] = useState<{
    query?: string | null
    startDate?: Date | null
    endDate?: Date | null
  }>({})



  // Create a stable filters key so that changing filters changes the query cache key
  const filtersKey = {
    query: appliedFilters.query ?? null,
    startDate: appliedFilters.startDate ? appliedFilters.startDate.toISOString() : null,
    endDate: appliedFilters.endDate ? appliedFilters.endDate.toISOString() : null,
  }

  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetching,
    isFetchingNextPage,
  } = useInfiniteQuery({
    initialPageParam: "",
    queryKey: ['logs', filtersKey],
    queryFn: async ({ pageParam = "" }: { pageParam?: string }) => {
      const params = new URLSearchParams()
      params.append('source', 'gateway-client')
      params.append('pageSize', '10')
      if (pageParam) params.append('bookmark', String(pageParam))
      if (filtersKey.query) params.append('query', String(filtersKey.query))
      if (filtersKey.startDate) params.append('startDate', String(filtersKey.startDate))
      if (filtersKey.endDate) params.append('endDate', String(filtersKey.endDate))

      const res = await fetch(`${API_BASE_URL}/log?` + params.toString())
      return res.json()
    },
    getNextPageParam: (lastPage: LogsResponse) => lastPage.hasNextPage ? lastPage.bookmark : null,
  })

  return (
    <div className="max-w-7xl">
      <LogFilters appliedFilters={appliedFilters} setAppliedFilters={setAppliedFilters} />
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-[500px]">Content</TableHead>
            <TableHead className="w-[100px]">Verified</TableHead>
            <TableHead className="w-[100px]">Client</TableHead>
            <TableHead className="w-[100px]">Timestamp</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {isFetching && !data?.pages?.length ? (
            <div>
              Loading...
            </div>
          ) : (
            data?.pages?.length === 0 ? (
              <div>
                No logs found
              </div>
            ) : (
              data?.pages?.flatMap(page => (page as LogsResponse).logs)?.map((log, index) => (
                <TableRow key={index} className={log.IsValid ? undefined : 'bg-red-100/60 hover:bg-red-200/70'}>
                  <TableCell>{log.Content}</TableCell>
                  <TableCell>{log.IsValid ? 'Yes' : 'No'}</TableCell>
                  <TableCell>{log.Source}</TableCell>
                  <TableCell>{new Date(log.Timestamp).toLocaleString()}</TableCell>
                </TableRow>
              ))
            )
          )}
        </TableBody>
      </Table>
      <div className="flex flex-row justify-between items-center mt-5">
        <div>
          <Button disabled={isFetchingNextPage || !hasNextPage} onClick={() => {
            fetchNextPage()
          }}>
            Load More
          </Button>
        </div>
        <div className="flex items-center gap-3">
          <div className="w-4 h-4 rounded-sm bg-red-200 border" aria-hidden="true" />
          <div className="text-sm text-muted-foreground">Red rows are invalid</div>
        </div>
      </div>
    </div>
  );
}

export default Logs;
