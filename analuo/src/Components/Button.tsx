import { styled } from "@mui/material";
import Button, { ButtonProps } from "@mui/material/Button";

export const PrimaryButton = styled(Button)<ButtonProps>(() => ({
  color: '#FFF8F0',
  backgroundColor: '#3BB273',
  '&:hover': {
    backgroundColor: '#339963',
  },
}));

export const SecondaryButton = styled(Button)<ButtonProps>(() => ({
  color: '#3BB273',
  backgroundColor: '#00000000',
  borderColor: '#3BB273',
  '&:hover': {
    borderColor: '#339963',
    backgroundColor: '#00000008',
  },
}));
