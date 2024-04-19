import styled from '@emotion/styled';
import { ChangeEventHandler } from 'react';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import { PrimaryButton } from './Button';

const VisuallyHiddenInput = styled('input')({
  clip: 'rect(0 0 0 0)',
  clipPath: 'inset(50%)',
  height: 1,
  overflow: 'hidden',
  position: 'absolute',
  bottom: 0,
  left: 0,
  whiteSpace: 'nowrap',
  width: 1,
});

export default function InputFileUpload({ onChange }: { onChange: ChangeEventHandler<HTMLInputElement> }) {

  return (
    <PrimaryButton
      component="label"
      role={undefined}
      variant="contained"
      tabIndex={-1}
      startIcon={<CloudUploadIcon />}
      sx={{ background: '3BB273' }}
    >
      Upload file
      <VisuallyHiddenInput type="file" onChange={onChange} />
    </PrimaryButton>
  );
}