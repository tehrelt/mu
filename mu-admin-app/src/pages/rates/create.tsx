import { RateCreateForm } from "@/components/forms/rate-create-form";
import {
  Breadcrumb,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { routes } from "@/shared/routes";
import { Link } from "react-router-dom";

export const RateCreatePage = () => {
  return (
    <>
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbLink asChild>
            <Link to={routes.rate.list}>Тарифы</Link>
          </BreadcrumbLink>
          <BreadcrumbSeparator />
          <BreadcrumbPage>Создать тариф</BreadcrumbPage>
        </BreadcrumbList>
      </Breadcrumb>
      <div className="flex">
        <RateCreateForm />
      </div>
    </>
  );
};
