import { useState } from 'react';

import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Container from '@mui/material/Container';
import Stack from '@mui/material/Stack';
import Stepper from '@mui/material/Stepper';
import Step from '@mui/material/Step';
import StepLabel from '@mui/material/StepLabel';
import Typography from '@mui/material/Typography';

import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import ErrorIcon from '@mui/icons-material/Error';
import CloseIcon from '@mui/icons-material/Close';

import { Address } from './Types/Address';
import { basePatterns, Pattern } from './Types/Pattern';

import { PrimaryButton, SecondaryButton } from './Components/Button';
import InputFileUpload from './Components/InputFileUpload';
import DisplayTable from './Components/DisplayTable'
import { TableSkeleton } from './Components/Skeletons';
import SettingsDialog from './Components/SettingsDialog';

import './App.css';
import Snackbar from '@mui/material/Snackbar';
import IconButton from '@mui/material/IconButton';
import Paper from '@mui/material/Paper';

const steps = [
  { label: 'Upload address file', cta: 'Validate'},
  { label: 'Confirm results', cta: 'Confirm'},
];

const successMessages = {
  icon: <CheckCircleIcon sx={{ fontSize: 128, color: '', opacity: .7 }} />,
  title: 'Hooray!',
  paragraph: 'Your changes have been whisked away into our system and are ready to be used.',
};

const failedMessage = {
  icon: <ErrorIcon sx={{ fontSize: 128 }} />,
  title: 'Whoops!',
  paragraph: 'It looks like an error occured during the update. Please make sure the your internet connections is stable and try again.',
};

const mapAddressRow = (a: Address) => ([
  a.Id,
  a.StreetNumber,
  a.StreetName,
  a.Unit,
  a.City,
  a.Zip,
  a.State,
  `${a.Region.Elementary} | ${a.Region.Middle} | ${a.Region.High}`,
  a.InCounty,
]);

async function parseCSV(file: File) {
  const headers: string[] = [];
  const dataRows: string[][] = [];

  const ab = await file.arrayBuffer();

  const decoder = new TextDecoder();
  const text = decoder.decode(ab);
          
  const fileRows = text.split('\r\n');
  headers.push(...fileRows[0].split(','));

  for (let i = 1; i < fileRows.length; i += 1) {
    const fields = [];
    const values = fileRows[i].split(',');

    for (let j = 0; j < values.length; j += 1) {
      const tFields = [ ];

      let hasEscapeStart = (index: number) => values[index].indexOf('"') === 0 && values[index].indexOf('"', values[index].indexOf('"') + 1) === -1;
      let hasEscapeEnd = (index: number) => values[index].indexOf('"') === values[index].length - 1;
      if (hasEscapeStart(j)) {
        do {
          tFields.push(values[j].replace('"', ''));
          j += 1;
        } while (!hasEscapeEnd(j));
      }

      tFields.push(values[j].replace('"', ''));
      fields.push(tFields.join(','));
    }
    dataRows.push(fields);
  }

  return { headers, rows: dataRows };
}

async function postData(apiUrl: string, data: Object) {
  try {
    const response = await fetch(
      apiUrl,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      }
    );

    if (!response.ok) {
      throw new Error(`Error fetching data: ${response.statusText}`);
    }

    const responseData = await response.json();
    return responseData;

  } catch (error) {
    console.log(error);
    return '[]';
  }
};

function App() {
  const [activeStep, setActiveStep] = useState(0);
  const [inputfile, setInputFile] = useState(new File([], ''));
  const [validation, setValidation] = useState({
    headers: [''],
    rows: [['']],
    url: '',
  });
  const [commit, setCommit] = useState({
    success: false,
    counts: {
      adds: 0,
      removes: 0,
    },
  });
  const [isLoading, setIsLoading] = useState(false);
  const [settingsOpen, setSettingsOpen] = useState(false);
  const [patterns, setPatterns] = useState(basePatterns);
  const [snackbar, setSnackbar] = useState({ open: false, message: '' });

  const canAdvance = () => activeStep !== 0 || inputfile.size !== 0;

  const handleNext = async () => {
    if (activeStep === 0) {
      setIsLoading(true);
      parseCSV(inputfile).then(async ({ headers, rows }) => {
        const body = {
          headers: headers,
          rows: rows,
        };
        console.log(`Patterns will be used in future updates: ${patterns}`)
        const response = await postData('http://localhost:3000/api/validations', body);
      
        try {
          const responseHeaders = Object.keys(response.data[0]);
          const responseRows = response.data.map((a: Address) => mapAddressRow(a));
  
          setValidation({
            headers: responseHeaders,
            rows: responseRows,
            url: response.url,
          });
        } catch (error) {
          setIsLoading(false);
        }
      
        setIsLoading(false);
      });
    }
    if (activeStep === 1) {
      setIsLoading(true);
      const body = {
        validationId: validation.url.replace('data/export/', '').replace('.csv', ''),
      };

      handleSnackbarNotification('Committing these changes');

      const response = await postData('http://localhost:3000/api/commits', body);

      setCommit({
        success: response.status === 200,
        counts: {
          adds: response.addCount,
          removes: response.removeCount
        }
      });
      setIsLoading(false);
    }

    setActiveStep((prevActiveStep) => prevActiveStep + 1);
  };

  const handleBack = () => {
    if (activeStep === 1) {
      setValidation({ headers: [], rows: [], url: '' });
    }
    setActiveStep((prevActiveStep) => prevActiveStep - 1);
  };

  const handleReset = () => {
    setActiveStep(0);
  };

  const handleFileChange = (event: any) => {
    setInputFile(event.target.files[0]);
    handleSnackbarNotification('Address file selected.')
  }

  const handleSettingsClick = () => {
    setSettingsOpen(true);
  }

  const handleSettingsClose = () => {
    setSettingsOpen(false);
  };

  const handleSettingsSave = (patterns: Pattern[]) => {
    setPatterns(patterns);
    setSettingsOpen(false);

    handleSnackbarNotification('Pattern settings saved.');
  }

  const handleSnackbarNotification = (message: string) => {
    setSnackbar({ open: true, message });
  }

  const handleSnackbarClose = () => setSnackbar({ open: false, message: '' });

  return (
    <Container maxWidth="md" sx={{ pt: 4 }}>
      <Paper variant="outlined" sx={{ p: 4, backgroundColor: '#FFF8F0' }}>
        <Stack>
          <Box sx={{ width: '100%' }}>
            <Stepper activeStep={activeStep}>
              {steps.map(({ label }) => {
                const stepProps: { completed?: boolean } = {};
                const labelProps: {
                  optional?: React.ReactNode;
                } = {};
                return (
                  <Step key={label} {...stepProps} sx={{ color: '#FFF8F0' }}>
                    <StepLabel {...labelProps}>{label}</StepLabel>
                  </Step>
                );
              })}
            </Stepper>
            {activeStep === steps.length ? (
              <Stack direction="column" justifyContent="space-between" alignContent="center" sx={{ p: 10, minHeight: 350 }}>
                <Box sx={{ m: 4, ml: 6, mr: 6 }} >
                  <Stack direction="row" justifyContent="center" sx={{ mb: 4}} >{commit.success ? successMessages.icon : failedMessage.icon}</Stack>
                  <Typography variant="h2" sx={{ fontWeight: 600, textAlign: 'center', mb: 2 }}>{commit.success ? successMessages.title : failedMessage.title}</Typography>
                  <Typography variant="body1" sx={{ textAlign: 'center' }}>{commit.success ? successMessages.paragraph : failedMessage.paragraph}</Typography>
                  {commit.success
                    ? (
                      <>
                        <Typography variant="body1" sx={{ textAlign: 'center' }}>Successfully added {commit.counts.adds} addresses.</Typography>
                        <Typography variant="body1" sx={{ textAlign: 'center' }}>Successfully removed {commit.counts.removes} addresses.</Typography>
                      </>
                    )
                  :''
                  }
                </Box>
                <Box sx={{ display: 'flex', flexDirection: 'row', pt: 2 }}>
                  <Box sx={{ flex: '1 1 auto' }} />
                  <SecondaryButton onClick={handleReset}>Reset</SecondaryButton>
                </Box>
              </Stack>
            ) : (
              <Stack flexDirection="column" justifyContent="space-between" sx={{ minHeight: 350 }} >
                <Typography sx={{ mt: 2, mb: 1 }}>Step {activeStep + 1}</Typography>
                <Stack flexDirection="row" justifyContent="center" >
                  {activeStep === 0
                    ? (
                      <Stack flexDirection="row" justifyContent="space-between" sx={{ width: "100%", maxWidth: 350 }}>
                        <SecondaryButton variant="outlined" onClick={handleSettingsClick}>
                          Adjust Settings
                        </SecondaryButton>
                        <InputFileUpload onChange={handleFileChange} />
                      </Stack>
                    )
                    : ''
                  }
                  {activeStep === 1 && validation.headers.length
                    ? (isLoading
                      ? <TableSkeleton />
                      : <DisplayTable
                        headers={validation.headers}
                        rows={validation.rows}
                        downloadUrl={validation.url}
                      />
                    )
                    : ''
                  }
                </Stack>
                <Stack direction='row' spacing={2} sx={{ pt: 2 }}>
                  <Button
                    color="inherit"
                    variant="outlined"
                    disabled={activeStep === 0}
                    onClick={handleBack}
                    sx={{ mr: 1 }}
                  >
                    Back
                  </Button>
                  <Box sx={{ flex: '1 1 auto' }} />
                  <PrimaryButton variant="contained" onClick={handleNext} disabled={!canAdvance()}>
                    {steps[activeStep].cta}
                  </PrimaryButton>
                </Stack>
              </Stack>
            )}
          </Box>
        </Stack>
      </Paper>
      {activeStep === 0
        ? (
          <SettingsDialog
            open={settingsOpen}
            handleClose={handleSettingsClose}
            handleSave={handleSettingsSave}
          />
        )
        : ''
      }
      <Snackbar
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
        open={snackbar.open}
        autoHideDuration={4000}
        onClose={handleSnackbarClose}
        message={snackbar.message}
        action={
          <IconButton
            size="small"
            aria-label="close"
            color="inherit"
            onClick={handleSnackbarClose}
          >
            <CloseIcon fontSize="small" />
          </IconButton>
        }
      />
    </Container>
  );
}

export default App;
