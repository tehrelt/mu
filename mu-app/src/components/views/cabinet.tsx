import { Cabinet } from "@/shared/types/cabinet";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "../ui/card";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { rateService } from "@/shared/services/rate.service";
import { Balance } from "../ui/balance";
import { Rate } from "@/shared/types/rate";
import { datef } from "@/shared/lib/utils";
import { InputOTP, InputOTPSlot } from "../ui/input-otp";
import React from "react";
import { Button } from "../ui/button";
import { ArrowDown, ArrowDownLeft, ArrowUp } from "lucide-react";
import { cabinetService } from "@/shared/services/cabinet.service";
import { toast } from "sonner";

type Props = {
  cabinet: Cabinet;
};

export const CabinetViewer = ({ cabinet }: Props) => {
  const query = useQuery({
    queryKey: [
      "cabinet",
      cabinet.id,
      { serviceId: cabinet.serviceId, accountId: cabinet.accountId },
    ],
    queryFn: async () => {
      const service = await rateService.find(cabinet.serviceId);
      return { service };
    },
  });

  return (
    <div className="flex gap-6">
      <div>
        <Card className="w-[380px]">
          <CardHeader>
            <CardTitle>Потреблено</CardTitle>
          </CardHeader>
          <CardContent className="flex">
            <div>
              {cabinet.consumed} {query.data?.service.measureUnit}
            </div>
          </CardContent>
        </Card>
      </div>
      {query.data && (
        <>
          <div>
            <Card className="w-[380px]">
              <CardHeader>
                <CardTitle>Поставщик {query.data.service.name}</CardTitle>
                <CardDescription></CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex gap-2">
                  <Balance balance={query.data.service.rate} />
                  <span>за {query.data.service.measureUnit}</span>
                </div>
              </CardContent>
            </Card>
          </div>
          <div>
            <ConsumptionRegister
              cabinet={cabinet}
              service={query.data.service}
            />
          </div>
        </>
      )}
    </div>
  );
};

const ConsumptionRegister = ({
  cabinet,
  service,
}: {
  cabinet: Cabinet;
  service: Rate;
}) => {
  const [value, setValue] = React.useState(cabinet.consumed);
  const queryClient = useQueryClient();

  const deltaValue = (delta: number) => {
    const res = value + delta;
    if (res < cabinet.consumed) return setValue(cabinet.consumed);

    setValue(res);
  };

  const { mutate: consume, isPending } = useMutation({
    mutationFn: async (val: { id: string; value: number }) => {
      await cabinetService.consume(val.id, val.value);
    },
    onSuccess: () => {
      toast.success("Показания отправлены");
    },
    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: ["cabinet"],
      });
    },
  });

  const handleSubmit = async () => {
    consume({ id: cabinet.id, value: value - cabinet.consumed });
  };

  const submitDisabled =
    value === cabinet.consumed || value < cabinet.consumed || isPending;

  return (
    <Card className="">
      <CardHeader>
        <CardTitle>Отправить показания</CardTitle>
      </CardHeader>
      <CardContent className="grid gap-2">
        <div className="grid grid-cols-8">
          {[...Array(8)].map((v, i) => (
            <Button
              className="row-1"
              variant={"ghost"}
              onClick={() => deltaValue(Math.pow(10, 8 - i - 1))}
            >
              <ArrowUp />
            </Button>
          ))}
          <div className="row-2 col-span-8 w-full">
            <InputOTP
              inputMode="numeric"
              maxLength={8}
              value={value.toString().padStart(8)}
              className="grid grid-cols-8"
            >
              {[...Array(8)].map((v, i) => (
                <InputOTPSlot index={i} className="" />
              ))}
            </InputOTP>
          </div>
          {[...Array(8)].map((v, i) => (
            <Button
              className="row-3"
              variant={"ghost"}
              onClick={() => deltaValue(-Math.pow(10, 8 - i - 1))}
            >
              <ArrowDown />
            </Button>
          ))}
        </div>

        <div className="flex">
          <Button
            className="w-full"
            disabled={submitDisabled}
            onClick={() => handleSubmit()}
          >
            {!submitDisabled ? (
              <span>
                Отправить {value - cabinet.consumed} {service.measureUnit}
              </span>
            ) : isPending ? (
              <span>Отправка...</span>
            ) : (
              <span>Введите новые показания</span>
            )}
          </Button>
        </div>
      </CardContent>
      <CardFooter>
        <CardDescription>
          {cabinet.updatedAt && (
            <p>Последнее показание: {datef(cabinet.updatedAt)}</p>
          )}
        </CardDescription>
      </CardFooter>
    </Card>
  );
};
