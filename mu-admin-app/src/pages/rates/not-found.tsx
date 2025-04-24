type Props = {
  id?: string;
};

export const RateNotFoundPage = ({ id }: Props) => {
  return (
    <div className="flex items-center justify-center gap-x-2 h-screen">
      <div className="flex flex-col items-center">
        <p className="text-8xl opacity-20 font-black ">404</p>
        <p className="text-xl font-bold ">Тариф не найден</p>
        {id && <p className="text-muted-foreground text-sm">ID тарифа: {id}</p>}
      </div>
    </div>
  );
};
