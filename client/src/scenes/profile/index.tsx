/* eslint-disable @typescript-eslint/no-explicit-any */
import { useEffect, useState } from "react";
import {
  Box,
  Stack,
  TextField,
  Typography,
  Button,
  Avatar,
  CircularProgress,
  Alert,
  Paper,
} from "@mui/material";
import Grid from "@mui/material/Grid";
import EditIcon from "@mui/icons-material/Edit";
import AccountCircleIcon from "@mui/icons-material/AccountCircle";
import { useFormik } from "formik";
import * as Yup from "yup";
import api from "../../api/axios"; 
import { btnstyle, activeTextFieldSx } from "../../styles";
import ErrorPage from "../../components/errorPages/errors";

interface ProfileData {
  full_name: string;
  date_of_birth: string;
  aadhaar_number: string;
  phone_number: string;
  address: string;
}

const validationSchema = Yup.object({
  full_name: Yup.string().min(3).max(20).required("Full name is required"),
  date_of_birth: Yup.date().required("Date of birth is required"),
  aadhaar_number: Yup.string().matches(/^\d{12}$/, "Aadhaar must be 12 digits").required(),
  phone_number: Yup.string().matches(/^[6-9]\d{9}$/, "Invalid phone number").required(),
  address: Yup.string().min(10).required(),
});

function Profile() {
  const [loading, setLoading] = useState(true);
  const [errorStatus, setErrorStatus] = useState<"404" | "500" | null>(null);
  const [isEditMode, setIsEditMode] = useState(false);
  const [hasExistingProfile, setHasExistingProfile] = useState(false);
  const [serverMsg, setServerMsg] = useState<{ type: "success" | "error"; text: string } | null>(null);

  const formik = useFormik<ProfileData>({
    initialValues: {
      full_name: "",
      date_of_birth: "",
      aadhaar_number: "",
      phone_number: "",
      address: "",
    },
    validationSchema: validationSchema,
    onSubmit: async (values, { setErrors }) => {
      try {
        setServerMsg(null);
        const payload = {
          ...values,
          date_of_birth: values.date_of_birth ? `${values.date_of_birth}T00:00:00Z` : ""
        };

        if (hasExistingProfile) {
          await api.put("restricted/profile", payload);
          setServerMsg({ type: "success", text: "Profile updated successfully!" });
        } else {
          await api.post("restricted/profile", payload);
          setServerMsg({ type: "success", text: "Profile created successfully!" });
          setHasExistingProfile(true);
        }
        setIsEditMode(false);
      } catch (err: any) {
        const backendErrors = err.response?.data;
        
        if (backendErrors && typeof backendErrors === 'object' && !backendErrors.error) {
          // 1. Map errors to individual fields (TextField red highlight)
          setErrors(backendErrors);
          
          // 2. Extract specific values from the map to show in the Alert box
          const errorList = Object.values(backendErrors).join(", ");
          setServerMsg({ type: "error", text: errorList });
        } else {
          setServerMsg({ type: "error", text: err.response?.data?.error || "Failed to save profile" });
        }
      }
    },
  });

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        setLoading(true);
        const res = await api.get("restricted/profile");
        if (res.data && Object.keys(res.data).length > 0) {
          const formattedData = {
            ...res.data,
            date_of_birth: res.data.date_of_birth ? res.data.date_of_birth.split("T")[0] : ""
          };
          formik.setValues(formattedData);
          setHasExistingProfile(true);
          setIsEditMode(false);
        } else {
          setHasExistingProfile(false);
          setIsEditMode(true);
        }
      } catch (err: any) {
        if (err.response?.status === 404) {
          setHasExistingProfile(false);
          setIsEditMode(true); 
        } else if (err.response?.status === 401) {
          window.location.href = "/login";
        } else {
          setErrorStatus("500");
        }
      } finally {
        setLoading(false);
      }
    };
    fetchProfile();
  }, []);

  if (errorStatus) return <ErrorPage code={errorStatus} />;
  if (loading) return (
    <Box display="flex" justifyContent="center" alignItems="center" height="100vh">
      <CircularProgress sx={{ color: "#E94057" }} />
    </Box>
  );

  return (
    <Box sx={{ minHeight: "100vh", bgcolor: "transparent", p: { xs: 2, md: 8 } }}>
      <Paper
        elevation={0}
        sx={{
          maxWidth: "1200px",
          margin: "0 auto",
          p: { xs: 4, md: 8 },
        }}
      >
        <Box sx={{ mb: 8, display: "flex", justifyContent: "space-between", alignItems: "flex-start" }}>
          <Stack direction="row" spacing={3} alignItems="center">
            <Avatar sx={{ bgcolor: "#FF6F61", width: 50, height: 50 }}>
              <AccountCircleIcon sx={{ fontSize: 50 }} />
            </Avatar>
            <Box>
              <Typography variant="h5" fontWeight={800} letterSpacing="-1px">
                {hasExistingProfile ? "My Profile" : "Create Account"}
              </Typography>
              <Typography variant="body1" color="text.secondary">
                {hasExistingProfile ? "Manage your identity and contact info" : "Complete your details to get started"}
              </Typography>
            </Box>
          </Stack>
          
          {!isEditMode && hasExistingProfile && (
            <Button
              variant="outlined"
              startIcon={<EditIcon />}
              onClick={() => setIsEditMode(true)}
              sx={{ borderRadius: "12px", color: "#E94057", borderColor: "#E94057", textTransform: "none", px: 3 }}
            >
              Edit Profile
            </Button>
          )}
        </Box>

        {serverMsg && (
          <Alert severity={serverMsg.type} sx={{ mb: 6, borderRadius: "16px" }}>
            {serverMsg.text}
          </Alert>
        )}

        <Box component="form" onSubmit={formik.handleSubmit}>
          <Grid container spacing={6}>
            <Grid size={{ xs: 12, md: 6 }}>
              <TextField
                fullWidth
                label="Full Name"
                variant="standard"
                slotProps={{ input: { readOnly: !isEditMode } }}
                sx={activeTextFieldSx}
                {...formik.getFieldProps("full_name")}
                error={formik.touched.full_name && Boolean(formik.errors.full_name)}
                helperText={formik.touched.full_name && formik.errors.full_name}
              />
            </Grid>

            <Grid size={{ xs: 12, md: 6 }}>
              <TextField
                fullWidth
                label="Date of Birth"
                type="date"
                variant="standard"
                InputLabelProps={{ shrink: true }}
                slotProps={{ input: { readOnly: !isEditMode } }}
                sx={activeTextFieldSx}
                {...formik.getFieldProps("date_of_birth")}
                error={formik.touched.date_of_birth && Boolean(formik.errors.date_of_birth)}
                helperText={formik.touched.date_of_birth && formik.errors.date_of_birth}
              />
            </Grid>

            <Grid size={{ xs: 12, md: 6 }}>
              <TextField
                fullWidth
                label="Aadhaar Number"
                variant="standard"
                slotProps={{ input: { readOnly: !isEditMode } }}
                sx={activeTextFieldSx}
                {...formik.getFieldProps("aadhaar_number")}
                error={formik.touched.aadhaar_number && Boolean(formik.errors.aadhaar_number)}
                helperText={formik.touched.aadhaar_number && formik.errors.aadhaar_number}
              />
            </Grid>

            <Grid size={{ xs: 12, md: 6 }}>
              <TextField
                fullWidth
                label="Phone Number"
                variant="standard"
                slotProps={{ input: { readOnly: !isEditMode } }}
                sx={activeTextFieldSx}
                {...formik.getFieldProps("phone_number")}
                error={formik.touched.phone_number && Boolean(formik.errors.phone_number)}
                helperText={formik.touched.phone_number && formik.errors.phone_number}
              />
            </Grid>

            <Grid size={{xs:12}}>
              <TextField
                fullWidth
                label="Residential Address"
                variant="standard"
                multiline
                rows={3}
                slotProps={{ input: { readOnly: !isEditMode } }}
                sx={activeTextFieldSx}
                {...formik.getFieldProps("address")}
                error={formik.touched.address && Boolean(formik.errors.address)}
                helperText={formik.touched.address && formik.errors.address}
              />
            </Grid>
          </Grid>

          {isEditMode && (
            <Box sx={{ mt: 10, display: "flex", justifyContent: "flex-end", gap: 3 }}>
              {hasExistingProfile && (
                <Button 
                  size="large"
                  onClick={() => { setIsEditMode(false); formik.resetForm(); }} 
                  sx={{ color: "#636e72", textTransform: "none", fontWeight: 600 }}
                >
                  Discard Changes
                </Button>
              )}
              <Button 
                type="submit" 
                variant="contained" 
                style={btnstyle} 
                sx={{ borderRadius: "100px", px: 10, py: 2, fontSize: "1.1rem" }} 
                disabled={formik.isSubmitting}
              >
                {formik.isSubmitting ? (
                  <CircularProgress size={24} color="inherit" />
                ) : (
                  hasExistingProfile ? "Update Profile" : "Create Profile"
                )}
              </Button>
            </Box>
          )}
        </Box>
      </Paper>
    </Box>
  );
}

export default Profile;