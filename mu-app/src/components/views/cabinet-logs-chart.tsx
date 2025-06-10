"use client";

import { Area, AreaChart, CartesianGrid, XAxis } from "recharts";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { datef } from "@/shared/lib/utils";
import { CabinetLog } from "@/shared/types/cabinet";

const chartConfig = {
  desktop: {
    label: "Desktop",
    color: "hsl(var(--chart-1))",
  },
  mobile: {
    label: "Mobile",
    color: "hsl(var(--chart-2))",
  },
} satisfies ChartConfig;

type Service = {
  name: string;
  logs: CabinetLog[];
};

type Props = {
  services: Service[];
};

type AggregatedLog = {
  [serviceName: string]: number;
} & { date: string | Date };

function aggregateLogsByDate(services: Service[]): AggregatedLog[] {
  const dateMap: Map<string, AggregatedLog> = new Map();

  for (const service of services) {
    for (const log of service.logs) {
      const dateKey = new Date(log.createdAt).toISOString().split("T")[0]; // только дата, без времени

      if (!dateMap.has(dateKey)) {
        //@ts-ignore
        dateMap.set(dateKey, { date: dateKey });
      }

      const aggregated = dateMap.get(dateKey)!;
      aggregated[service.name] = (aggregated[service.name] || 0) + log.consumed;
    }
  }

  return Array.from(dateMap.values()).sort((a, b) =>
    //@ts-ignore
    a.date.localeCompare(b.date)
  );
}

export function LogsChart({ services }: Props) {
  const data = aggregateLogsByDate(services);

  console.log(data);

  return (
    <Card>
      <CardHeader>
        <CardTitle>
          Потребление по дням ({services.map((s) => s.name).join(", ")})
        </CardTitle>
        <CardDescription>
          Показывает потребление по дням за последний месяц
        </CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <AreaChart
            accessibilityLayer
            data={data}
            margin={{
              left: 12,
              right: 12,
            }}
          >
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="date"
              tickLine={true}
              axisLine={false}
              tickMargin={8}
              tickFormatter={(value) => datef(value, "DD/MM/YY")}
            />
            <ChartTooltip
              cursor={false}
              content={<ChartTooltipContent indicator="dot" />}
            />
            {services.map((s) => (
              <Area
                dataKey={`${s.name}`}
                type="natural"
                fill="var(--color-mobile)"
                fillOpacity={0.4}
                stroke="var(--color-mobile)"
                stackId={s.name}
              />
            ))}
          </AreaChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
