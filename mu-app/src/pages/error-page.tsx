import { useRouteError } from "react-router-dom";

export default function ErrorPage() {
  const error = useRouteError() as Error;
  console.error(error);

  return (
    <div id="error-page">
      <h1>Oops!</h1>
      <p>Sorry, an unexpected error has occurred.</p>
      <p>
        <p>{error.name}</p>
        <p>{error.message}</p>
        <p>{error.stack}</p>
      </p>
    </div>
  );
}
