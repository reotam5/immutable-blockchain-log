import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import type { Theme } from "@emotion/react"
import { TextField } from '@mui/material'
import type { SxProps } from "@mui/material/styles"
import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker'
import type { PickerValue } from "@mui/x-date-pickers/internals"
import { X } from "lucide-react"
import { useState } from "react"
import { Label } from "./ui/label"
import { Separator } from "./ui/separator"

export default function LogFilters() {
  const [query, setQuery] = useState("")
  const [startDate, setStartDate] = useState<PickerValue>(null)
  const [endDate, setEndDate] = useState<PickerValue>(null)


  const [appliedFilters, setAppliedFilters] = useState<{
    query?: string | null
    startDate?: Date | null
    endDate?: Date | null
  }>({})

  const style: SxProps<Theme> = {
    '.MuiPickersTextField-root': {
      height: '36px',
    },
    '& .MuiPickersInputBase-root': {
      height: '36px',
    },
    '& .MuiFormLabel-root': {
      top: '-6px',
    }
  }

  return (
    <div className="mb-6 ">
      <Label className="mb-3">Filters</Label>
      <div className="flex flex-col md:flex-row md:items-end gap-4">
        <div className="flex-1">
          <TextField
            id="search"
            variant="outlined"
            label="Search logs..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            sx={{
              width: '100%',
              '& .MuiInputBase-root': {
                height: '36px',
              },
              '& .MuiFormLabel-root': {
                top: '-6px',
              }
            }}
          />
        </div>

        <div className="flex gap-4">
          <div>
            <DateTimePicker
              label="Start date & time"
              value={startDate}
              onChange={(newValue) => setStartDate(newValue)}
              sx={style}
            />
          </div>

          <div>
            <DateTimePicker
              label="End date & time"
              value={endDate}
              onChange={(newValue) => setEndDate(newValue)}
              sx={style}
            />
          </div>
        </div>

        <div className="flex items-center gap-2">
          <Button
            onClick={() => {
              const next = {
                query: query || null,
                startDate: startDate?.toDate() || null,
                endDate: endDate?.toDate() || null,
              }
              setAppliedFilters(next)
            }}
          >
            Apply
          </Button>
          <Button
            variant="outline"
            onClick={() => {
              setQuery("")
              setStartDate(null)
              setEndDate(null)
              setAppliedFilters({})
            }}
          >
            Clear
          </Button>
        </div>
      </div>
      <div className="mt-3 flex flex-wrap gap-2 pb-2">
        {appliedFilters.query ? (
          <Badge>
            <span>{appliedFilters.query}</span>
            <Button
              variant="ghost"
              onClick={() => {
                setAppliedFilters((s) => ({ ...s, query: null }))
                setQuery("")
              }}
              className="p-0 m-0 h-5 w-5"
            >
              <X className="size-3" />
            </Button>
          </Badge>
        ) : null}

        {appliedFilters.startDate ? (
          <Badge>
            <span>Start date: {appliedFilters.startDate?.toLocaleString()}</span>
            <Button
              variant="ghost"
              onClick={() => {
                setAppliedFilters((s) => ({ ...s, startDate: null }))
                setStartDate(null)
              }}
              className="p-0 m-0 h-5 w-5"
            >
              <X className="size-3" />
            </Button>
          </Badge>
        ) : null}

        {appliedFilters.endDate ? (
          <Badge>
            <span>Etart date: {appliedFilters.endDate?.toLocaleString()}</span>
            <Button
              variant="ghost"
              onClick={() => {
                setAppliedFilters((s) => ({ ...s, endDate: null }))
                setEndDate(null)
              }}
              className="p-0 m-0 h-5 w-5"
            >
              <X className="size-3" />
            </Button>
          </Badge>
        ) : null}
      </div>
      <Separator className="my-1" />
    </div>
  )
}
