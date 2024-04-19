import { Skeleton, Stack } from "@mui/material";

export function TableSkeleton() {
  return (
    <Stack flexDirection="column" sx={{ width: "80%", mb: 4 }} >
      <Skeleton height={60} width="100%" />
      <Skeleton height={60} width="100%" />
      <Skeleton height={60} width="100%" />
      <Skeleton height={60} width="100%" />
      <Skeleton height={60} width="100%" />
      <Skeleton height={60} width="100%" />
      <Skeleton height={60} width="100%" />
    </Stack>
  );
}