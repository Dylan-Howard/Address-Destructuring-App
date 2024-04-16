import { useState } from 'react';
import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import Stack from '@mui/material/Stack';
import Stepper from '@mui/material/Stepper';
import Step from '@mui/material/Step';
import StepLabel from '@mui/material/StepLabel';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
// import logo from './logo.svg';
import './App.css';
import InputFileUpload from './Common/InputFileUpload';
import DisplayTable from './Common/DisplayTable'
import { Address } from './Common/Address';
import CircularProgress from '@mui/material/CircularProgress';

const steps = ['Upload address file', 'Confirm results', 'Upload changes'];

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
  const [file, setFile] = useState(new File([], ''));
  const [headers, setHeaders] = useState(['']);
  const [rows, setRows] = useState([['']]);
  const [validation, setValidation] = useState({
    headers: [''],
    rows: [['']],
  });
  const [uploadStatus, setUploadStatus] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const canAdvance = () => activeStep !== 0 || file.size !== 0;

  const handleNext = async () => {
    if (activeStep === 0) {
      parseCSV(file)
        .then(({ headers, rows }) => {
          setHeaders(headers);
          setRows(rows);
        });
    }
    if (activeStep === 1) {
      setIsLoading(true);
      const body = {
        headers,
        rows,
      };

      const response = await postData('http://localhost:3000/api/addresses/validate', body);
      
      const responseHeaders = Object.keys(response.data[0]);
      const responseRows = response.data.map((a: Address) => ([
        a.Id,
        a.StreetNumber,
        a.StreetName,
        a.Unit,
        a.City,
        a.Zip,
        a.State,
        `${a.Region.Elementary} | ${a.Region.Middle} | ${a.Region.High}`,
        a.InCounty,
      ]));

      setValidation({
        headers: responseHeaders,
        rows: responseRows,
      });
      setIsLoading(false);
    }
    if (activeStep === 2) {
      setIsLoading(true);
      const body = {
        headers,
        rows: validation.rows,
      };
      const response = await postData('http://localhost:8080/api/addresses/submit', body);

      setUploadStatus(!!response);
      setIsLoading(false);
    }

    setActiveStep((prevActiveStep) => prevActiveStep + 1);
  };

  const handleBack = () => {
    if (activeStep === 1) {
      setRows([]);
    }
    if (activeStep === 2) {
      setValidation({ headers: [], rows: [] });
    }
    setActiveStep((prevActiveStep) => prevActiveStep - 1);
  };

  const handleReset = () => {
    setActiveStep(0);
  };

  const handleFileChange = (event: any) => {
    setFile(event.target.files[0]);
  }

  return (
    <Container maxWidth="md" sx={{ mt: 4 }}>
      <Stack>
        <Box sx={{ width: '100%' }}>
          <Stepper activeStep={activeStep}>
            {steps.map((label, index) => {
              const stepProps: { completed?: boolean } = {};
              const labelProps: {
                optional?: React.ReactNode;
              } = {};
              return (
                <Step key={label} {...stepProps}>
                  <StepLabel {...labelProps}>{label}</StepLabel>
                </Step>
              );
            })}
          </Stepper>
          {activeStep === steps.length ? (
            <Stack direction="column" justifyContent="center" alignContent="center">
              <Typography variant="h2" sx={{ fontWeight: 600 }}>{uploadStatus ? "Success" : "Whoops"}</Typography>
              <Box sx={{ display: 'flex', flexDirection: 'row', pt: 2 }}>
                <Box sx={{ flex: '1 1 auto' }} />
                <Button onClick={handleReset}>Reset</Button>
              </Box>
            </Stack>
          ) : (
            <>
              <Typography sx={{ mt: 2, mb: 1 }}>Step {activeStep + 1}</Typography>
              <Stack flexDirection="row" justifyContent="center">
                {activeStep === 0
                  ? <InputFileUpload onChange={handleFileChange} />
                  : ''
                }
                {activeStep === 1
                  ? (isLoading ? <CircularProgress /> : <DisplayTable headers={headers} rows={rows} />)
                  : ''
                }
                {activeStep === 2
                  ? (isLoading ? <CircularProgress /> : <DisplayTable headers={validation.headers} rows={validation.rows} />)
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
                <Button variant="contained" onClick={handleNext} disabled={!canAdvance()}>
                  {activeStep === steps.length - 1 ? 'Finish' : 'Next'}
                </Button>
              </Stack>
            </>
          )}
        </Box>
      </Stack>
    </Container>
  );
}

export default App;
