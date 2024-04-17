import { FormEvent, useState } from 'react';
import Button from '@mui/material/Button';
// import TextField from '@mui/material/TextField';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
// import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';

import patterns from '../data/patterns.json';
import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';
import Select, { SelectChangeEvent } from '@mui/material/Select';
import MenuItem from '@mui/material/MenuItem';

function patternToRow(
  pattern: {
    expression: string;
    type: string;
    map?: {
      delimiter?: number | undefined;
      descriptor?: number | undefined;
      startValue?: number | undefined;
      endValue?: number | undefined;
      unitValue?: number | undefined;
    } | undefined;
  },
  setPatterns: Function,
) {
  const handleChange = (event: SelectChangeEvent) => {
    setPatterns(event.target.value as string);
  };

  return (
    <Box key={pattern.expression} sx={{ p:1, pl:2, pr:2, mb:2 }}>
      <Stack>
        <Typography sx={{ fontSize: 12, opacity:.7 }}>Expression</Typography>
        <Typography sx={{ mb: 2 }}>{pattern.expression}</Typography>
        <Stack flexDirection="row" justifyContent="space-between">
          <Box>
            <Typography sx={{ fontSize: 12, opacity:.7 }}>Type</Typography>
            <Typography>{pattern.type}</Typography>
          </Box>
          {pattern.type === 'list'
            ? (
              <FormControl size="small" sx={{ minWidth: 100 }}>
                <InputLabel id="select-delimiter-label">Delimiter</InputLabel>
                <Select
                  labelId="select-delimiter-label"
                  id="select-delimiter"
                  value={pattern.map?.delimiter ? pattern.map?.delimiter.toString() : undefined}
                  label="Delimiter"
                  onChange={handleChange}
                >
                  <MenuItem value={2}>2</MenuItem>
                  <MenuItem value={3}>3</MenuItem>
                  <MenuItem value={4}>4</MenuItem>
                </Select>
              </FormControl>
            )
            : ''
          }
          {pattern.type === 'unit'
            ? (
              <FormControl size="small" sx={{ minWidth: 100 }}>
                <InputLabel id="select-delimiter-label">Descriptor</InputLabel>
                <Select
                  labelId="select-delimiter-label"
                  id="select-delimiter"
                  value={pattern.map?.descriptor ? pattern.map?.descriptor.toString() : undefined}
                  label="Delimiter"
                  onChange={handleChange}
                >
                  <MenuItem value={2}>2</MenuItem>
                  <MenuItem value={3}>3</MenuItem>
                  <MenuItem value={4}>4</MenuItem>
                </Select>
              </FormControl>
            )
            : ''
          }
          {pattern.type === 'unit'
            ? (
              <FormControl size="small" sx={{ minWidth: 100 }}>
                <InputLabel id="select-delimiter-label">Unit</InputLabel>
                <Select
                  labelId="select-delimiter-label"
                  id="select-delimiter"
                  value={pattern.map?.unitValue ? pattern.map?.unitValue.toString() : undefined}
                  label="Delimiter"
                  onChange={handleChange}
                >
                  <MenuItem value={2}>2</MenuItem>
                  <MenuItem value={3}>3</MenuItem>
                  <MenuItem value={4}>4</MenuItem>
                </Select>
              </FormControl>
            )
            : ''
          }
          {pattern.type === 'range'
            ? (
              <FormControl size="small" sx={{ minWidth: 100 }}>
                <InputLabel id="select-delimiter-label">Unit</InputLabel>
                <Select
                  labelId="select-delimiter-label"
                  id="select-delimiter"
                  value={pattern.map?.startValue ? pattern.map?.startValue.toString() : undefined}
                  label="Delimiter"
                  onChange={handleChange}
                >
                  <MenuItem value={2}>2</MenuItem>
                  <MenuItem value={3}>3</MenuItem>
                  <MenuItem value={4}>4</MenuItem>
                </Select>
              </FormControl>
            )
            : ''
          }
          {pattern.type === 'range' && pattern.map?.endValue
            ? (
              <FormControl size="small" sx={{ minWidth: 100 }}>
                <InputLabel id="select-delimiter-label">Unit</InputLabel>
                <Select
                  labelId="select-delimiter-label"
                  id="select-delimiter"
                  value={pattern.map.endValue.toString()}
                  label="Delimiter"
                  onChange={handleChange}
                >
                  <MenuItem value={2}>2</MenuItem>
                  <MenuItem value={3}>3</MenuItem>
                  <MenuItem value={4}>4</MenuItem>
                </Select>
              </FormControl>
            )
            : ''
          }
          <Box>
            <Typography sx={{ fontSize: 12, opacity:.7 }}>Map</Typography>
            {/* <Typography>{pattern.map ? pattern.map.toString() : 'N/A'}</Typography>   */}
          </Box>
        </Stack>
      </Stack>
    </Box>
  )
}

export default function FormDialog({ open, handleClose}: { open: boolean, handleClose: Function }) {
  const [patternSettings, setPatternSettings] = useState(patterns.patterns);

  return (
    <Dialog
        open={open}
        onClose={() => handleClose()}
        PaperProps={{
          component: 'form',
          onSubmit: (event: FormEvent<HTMLFormElement>) => {
            event.preventDefault();
            const formData = new FormData(event.currentTarget);
            const formJson = Object.fromEntries((formData as any).entries());
            const email = formJson.email;
            console.log(email);
            handleClose();
          },
        }}
      >
        <DialogTitle>Patterns</DialogTitle>
        <DialogContent>
          {/* <DialogContentText>
            To subscribe to this website, please enter your email address here. We
            will send updates occasionally.
          </DialogContentText> */}
          {/* {patterns.patterns.map((pat) => (
            <Box key={pat.expression} sx={{ p:1, pl:2, pr:2, mb:2 }}>
              <Stack>
                <Typography sx={{ fontSize: 12, opacity:.7 }}>Expression</Typography>
                <Typography sx={{ mb: 1 }}>{pat.expression}</Typography>
                <Stack flexDirection="row" justifyContent="space-between">
                  <Box>
                    <Typography sx={{ fontSize: 12, opacity:.7 }}>Type</Typography>
                    <Typography>{pat.type}</Typography>
                  </Box>
                  <Box>
                    <Typography sx={{ fontSize: 12, opacity:.7 }}>Map</Typography>
                    <Typography>{pat.map ? pat.map.toString() : 'N/A'}</Typography>  
                  </Box>
                </Stack>
              </Stack>
            </Box>
          ))} */}
          {patterns.patterns.map((pat) => patternToRow(pat, setPatternSettings))

          }
        </DialogContent>
        <DialogActions>
          <Button onClick={() => handleClose()}>Cancel</Button>
          <Button type="submit">Save</Button>
        </DialogActions>
      </Dialog>
  );
}