import ScenariosView from "./ScenariosView";

export default function ScenariosPage({
  searchParams,
}: {
  searchParams: { highlight?: string };
}) {
  return <ScenariosView highlight={searchParams.highlight} />;
}
