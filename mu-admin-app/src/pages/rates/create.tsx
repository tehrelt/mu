import { RateCreateForm } from "@/components/forms/rate-create-form";

export const RateCreatePage = () => {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen ">
      <div className=" p-8 rounded-lg shadow-md">
        <RateCreateForm />
      </div>
    </div>
  );
};
