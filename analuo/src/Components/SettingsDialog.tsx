import { ChangeEvent, useState } from 'react';

import { Pattern, basePatterns } from '../Types/Pattern';

import Box from '@mui/material/Box';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import Select, { SelectChangeEvent } from '@mui/material/Select';
import Stack from '@mui/material/Stack';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';
import DeleteIcon from '@mui/icons-material/Delete';
import IconButton from '@mui/material/IconButton';
import { PrimaryButton, SecondaryButton } from './Button';

function PatternRow( id: number, pattern: Pattern, handleUpdate: Function, handleRemove: Function ) {

  const handleSelectChange = (event: SelectChangeEvent) => {
    if (pattern.map) {
      const mapValues = Object.values(pattern.map);

      for (let i = 0; i < mapValues.length; i += 1) {
        if (mapValues[i].toString() === event.target.value.toString()) {
          event.preventDefault();
          console.log('Preventing duplicate capture group usage');
          return;
        }
      }
    }

    if (event.target.value) {
      const newMap = pattern.map || {};
      const targetField = event.target.name.replace(`pattern-${id}-`, '');

      if (targetField === 'unit') {
        newMap.unitValue = parseInt(event.target.value, 10);
      }
      if (targetField === 'delimiter') {
        newMap.delimiter = parseInt(event.target.value, 10);
      }
      if (targetField === 'descriptor') {
        newMap.descriptor = parseInt(event.target.value, 10);
      }
      if (targetField === 'startValue') {
        newMap.startValue = parseInt(event.target.value, 10);
      }
      if (targetField === 'endValue') {
        newMap.endValue = parseInt(event.target.value, 10);
      }

      handleUpdate( id, { ...pattern, map: newMap } );
    } else {
      handleUpdate( id, pattern );
    }
  };

  const handleTextChange = (event: ChangeEvent<HTMLInputElement>) => {
    handleUpdate(
      id,
      {
        expression: event.target.value,
        type: pattern.type,
        map: pattern.map,
      },
    );
  }

  return (
    <Box key={`pattern-${id}`} sx={{ pt:1, pb: 3, ml:2, mr:2, mb:3, borderBottom: 1, borderColor: '#E8E4F1' }}>
      <Stack>
        <Stack flexDirection="row" justifyContent="end" sx={{ mb:2 }} >
          <IconButton aria-label="remove" onClick={() => handleRemove(id)}>
            <DeleteIcon />
          </IconButton>
        </Stack>
        <TextField
          label="Expression"
          variant="outlined"
          size="small"
          value={pattern.expression}
          onChange={handleTextChange}
          sx={{ mb: 2 }}
        />
        <Stack flexDirection="row" justifyContent="space-between">
          <Box>
            <FormControl size="small" sx={{ minWidth: 100 }}>
                <InputLabel id="select-type-label">Type</InputLabel>
                <Select
                  labelId="select-type-label"
                  id="select-type"
                  value={pattern.type}
                  label="Type"
                  onChange={handleSelectChange}
                  size="small"
                >
                  <MenuItem value="unit">Unit</MenuItem>
                  <MenuItem value="range">Range</MenuItem>
                  <MenuItem value="list">List</MenuItem>
                  <MenuItem value="nested range">Nested Range</MenuItem>
                </Select>
              </FormControl>
          </Box>
          {pattern.type === 'list'
            ? (
              <FormControl size="small" sx={{ minWidth: 100 }}>
                <InputLabel id="select-delimiter-label">Delimiter</InputLabel>
                <Select
                  labelId="select-delimiter-label"
                  id="select-delimiter"
                  name={`pattern-${id}-delimiter`}
                  value={pattern.map?.delimiter ? pattern.map?.delimiter.toString() : '' }
                  label="Delimiter"
                  onChange={handleSelectChange}
                  size="small"
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
                  name={`pattern-${id}-descriptor`}
                  value={pattern.map?.descriptor ? pattern.map?.descriptor.toString() : '' }
                  label="Descriptor"
                  onChange={handleSelectChange}
                  size="small"
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
                  labelId="select-unit-label"
                  id="select-unit"
                  name={`pattern-${id}-unit`}
                  value={pattern.map?.unitValue ? pattern.map?.unitValue.toString() : '' }
                  label="Unit"
                  onChange={handleSelectChange}
                  size="small"
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
                <InputLabel id="select-startValue-label">Start Value</InputLabel>
                <Select
                  labelId="select-start-value-label"
                  id="select-start-value"
                  name={`pattern-${id}-start-value`}
                  value={pattern.map?.startValue ? pattern.map?.startValue.toString() : '' }
                  label="Start Value"
                  onChange={handleSelectChange}
                  size="small"
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
                <InputLabel id="select-end-value-label">End Value</InputLabel>
                <Select
                  labelId="select-end-value-label"
                  id="select-end-value"
                  name={`pattern-${id}-end-value`}
                  value={pattern.map.endValue.toString()}
                  label="End Value"
                  onChange={handleSelectChange}
                  size="small"
                >
                  <MenuItem value={2}>2</MenuItem>
                  <MenuItem value={3}>3</MenuItem>
                  <MenuItem value={4}>4</MenuItem>
                </Select>
              </FormControl>
            )
            : ''
          }
        </Stack>
      </Stack>
    </Box>
  )
}

export default function SettingsDialog({ open, handleClose, handleSave }: { open: boolean, handleClose: Function, handleSave: Function }) {
  const [patternSettings, setPatternSettings] = useState(basePatterns);

  const handleUpdate = (targetId: number, updatedPattern: Pattern) => {
    const newPatterns = patternSettings.map((pat, i) => {
      if (i === targetId) {
        return {
          ...updatedPattern,
          map: updatedPattern.map || {
            delimiter: 0,
            descriptor: 0,
            startValue: 0,
            endValue: 0,
            unitValue: 0,
          }
        };
      }
      return pat;
    });
    setPatternSettings(newPatterns);
  };

  const handleDelete = (targetId: number) => {
    console.log(`deleting ${targetId}`);
    setPatternSettings(patternSettings.filter((pat, i) => i !== targetId));
  };

  return (
    <Dialog
        fullWidth
        open={open}
        onClose={() => handleClose()}
        maxWidth="md"
        sx={{ backgroundColor: '#FFF8F0' }}
      >
        <DialogTitle>Patterns</DialogTitle>
        <DialogContent>
          <Typography>The destructurer implements the following patterns in reverse order to parse the address data.</Typography>
          { patternSettings.map((pat, i) => PatternRow(i, pat, handleUpdate, handleDelete)) }
        </DialogContent>
        <DialogActions>
          <SecondaryButton onClick={() => handleClose()}>Cancel</SecondaryButton>
          <PrimaryButton type="submit" onClick={() => handleSave()}>Save</PrimaryButton>
        </DialogActions>
      </Dialog>
  );
}