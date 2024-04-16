import { ChangeEventHandler } from 'react';

export default function InputFileUpload({ onChange }: { onChange: ChangeEventHandler<HTMLInputElement> }) {

  return (
    <input type="file" name="addressFile" onChange={onChange} />
  );
}