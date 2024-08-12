import { useEffect, useState } from "react";
import "./App.css";
import { QRCode } from "react-qrcode-logo";
import { socket } from "./websocket";

function App() {
  const [qrCode, setQrCode] = useState("");

  useEffect(() => {
    socket.addEventListener("message", (event) => {
      console.log("message from server: ", event.data);
      setQrCode(event.data);
    });
  }, []); // eslint-disable-line

  return (
    <>
      <h1>Ini Ceritanya Buat Absen</h1>
      <div>
        <QRCode size={400} ecLevel="M" quietZone={20} value={qrCode} />
      </div>
      <p className="read-the-docs">Scan QR Above</p>
    </>
  );
}

export default App;
