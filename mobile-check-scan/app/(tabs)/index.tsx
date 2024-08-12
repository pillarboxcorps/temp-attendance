import {
  BarcodeScanningResult,
  CameraView,
  useCameraPermissions,
} from "expo-camera";
import { useFocusEffect } from "expo-router";
import { useCallback, useState } from "react";
import {
  StyleSheet,
  ScrollView,
  View,
  Text,
  Button,
  TouchableOpacity,
  Dimensions,
  Alert,
} from "react-native";

import { SafeAreaView } from "react-native-safe-area-context";

export default function HomeScreen() {
  const [permission, requestPermission] = useCameraPermissions();
  const [showCamera, setShowCamera] = useState(false);

  const handleBarcode = async (result: BarcodeScanningResult) => {
    const url = `/send-message?message=${result.data}`;
    try {
      const response = await fetch(url);
      if (!response.ok) Alert.alert(`Response status: ${response.status}`);

      // const parseRes = await response.json();
      //
      // Alert.alert("Mantap!!", parseRes);
    } catch (err: any) {
      Alert.alert("err: ", err.message);
    } finally {
      setShowCamera(false);
    }
  };

  useFocusEffect(
    useCallback(() => {
      return () => {
        setShowCamera(false);
      };
    }, []),
  ); // eslint-disable-line

  if (!permission) {
    // Camera permissions are still loading.
    return <View />;
  }

  if (!permission.granted) {
    // Camera permissions are not granted yet.
    return (
      <View style={styles.container}>
        <Text style={styles.message}>
          We need your permission to show the camera
        </Text>
        <Button onPress={requestPermission} title="grant permission" />
      </View>
    );
  }

  return (
    <SafeAreaView style={styles.safeAreaContainer}>
      <ScrollView>
        <View style={styles.qrScreen}>
          <View style={styles.topContainer}>
            <Text style={{ fontSize: 24, fontWeight: 700, color: "#474747" }}>
              Please Scan the QR Code
            </Text>
            <Text style={{ color: "#6e6e6e", marginTop: 4 }}>
              Scan the qr to process your doing
            </Text>
          </View>
          {showCamera ? (
            <View style={styles.container}>
              <CameraView
                style={styles.camera}
                barcodeScannerSettings={{
                  barcodeTypes: ["qr"],
                }}
                onBarcodeScanned={handleBarcode}
              />
            </View>
          ) : (
            <View />
          )}
          <View
            style={{
              display: "flex",
              justifyContent: "center",
              alignItems: "center",
            }}
          >
            <Text>
              <TouchableOpacity
                style={styles.button}
                onPress={() => setShowCamera(!showCamera)}
              >
                <Text style={styles.textBottom}>
                  {showCamera ? "Close" : "Open"} Camera
                </Text>
              </TouchableOpacity>
            </Text>
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  qrScreen: {
    padding: 12,
    height: Dimensions.get("window").height - 200,
    display: "flex",
    justifyContent: "center",
    gap: 14,
  },
  topContainer: {
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
  },
  safeAreaContainer: {
    backgroundColor: "#ffffff",
    height: "100%",
  },
  container: {
    justifyContent: "center",
    padding: 20,
  },
  message: {
    textAlign: "center",
    paddingBottom: 10,
  },
  camera: {
    borderRadius: 12,
    overflow: "hidden",
    height: 300,
  },
  button: {
    paddingVertical: 8,
    paddingHorizontal: 12,
    borderRadius: 10,
    borderWidth: 1,
    borderColor: "#000",
  },
  textBottom: {
    fontSize: 12,
    fontWeight: "bold",
    color: "black",
  },
});
