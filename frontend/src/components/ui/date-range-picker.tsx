import { useState } from "react"
import { format, startOfMonth, endOfMonth, subMonths } from "date-fns"
import { id } from "date-fns/locale"
import { CalendarIcon } from "lucide-react"
import type { DateRange } from "react-day-picker"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Calendar } from "@/components/ui/calendar"
import {
  Popover, PopoverContent, PopoverTrigger,
} from "@/components/ui/popover"

interface DateRangePickerProps {
  value?: DateRange
  onChange?: (range: DateRange | undefined) => void
  placeholder?: string
}

function getMonthPresets() {
  const now = new Date()
  const presets: { label: string; range: DateRange }[] = []
  for (let i = 0; i < 6; i++) {
    const date = subMonths(now, i)
    const from = startOfMonth(date)
    const to = i === 0 ? now : endOfMonth(date)
    presets.push({
      label: format(date, "MMMM yyyy", { locale: id }),
      range: { from, to },
    })
  }
  return presets
}

export function DateRangePicker({ value, onChange, placeholder = "Pilih periode" }: DateRangePickerProps) {
  const [open, setOpen] = useState(false)
  const presets = getMonthPresets()

  const handlePreset = (range: DateRange) => {
    onChange?.(range)
    setOpen(false)
  }

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          className={cn(
            "justify-start text-left font-normal",
            !value?.from && "text-muted-foreground"
          )}
        >
          <CalendarIcon className="mr-2 h-4 w-4" />
          {value?.from ? (
            value.to
              ? `${format(value.from, "dd MMM yyyy", { locale: id })} – ${format(value.to, "dd MMM yyyy", { locale: id })}`
              : format(value.from, "dd MMM yyyy", { locale: id })
          ) : placeholder}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-auto p-0" align="start">
        <div className="flex">
          <div className="flex flex-col gap-1 border-r px-3 py-3">
            {presets.map((preset) => (
              <button
                type="button"
                key={preset.label}
                className="rounded-md px-3 py-1.5 text-left text-sm capitalize hover:bg-accent hover:text-accent-foreground transition-colors"
                onClick={() => handlePreset(preset.range)}
              >
                {preset.label}
              </button>
            ))}
          </div>
          <div className="p-1">
            <Calendar
              mode="range"
              selected={value}
              onSelect={onChange}
              numberOfMonths={1}
            />
          </div>
        </div>
      </PopoverContent>
    </Popover>
  )
}
