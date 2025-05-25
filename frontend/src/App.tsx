import { useState } from "react";
import logo from "./assets/images/logo-universal.png";
import "./App.css";
import "./test.css";
import { Greet } from "../wailsjs/go/main/App";
import { CreateOrUpdate } from '../wailsjs/go/user/UserService';

function App() {
  const [resultText, setResultText] = useState(
    "Please enter your name below 👇"
  );
  const [name, setName] = useState("");
  const updateName = (e: any) => setName(e.target.value);
  const updateResultText = (result: string) => setResultText(result);

  function greet() {
    Greet(name).then(updateResultText);
  }

  const checkDb = async () => {
    const userId = await CreateOrUpdate({
      firstName: "John",
      lastName: "Doe",
    });

    console.log("User ID:", userId);
  };

  return (
    <div id="App">
      <img src={logo} id="logo" alt="logo" />
      <div className="p-8 bg-blue-500 text-red-500 font-bold test-custom">
        Tailwind works!
      </div>
      <div id="result" className="result">
        {resultText}
      </div>
      <div id="input" className="input-box">
        <input
          id="name"
          className="input"
          onChange={updateName}
          autoComplete="off"
          name="input"
          type="text"
        />
        <button className="btn" onClick={greet}>
          Greet
        </button>
      </div>

      <button
        onClick={checkDb}
        className="mt-2 px-4 py-2 bg-blue-500 text-white rounded"
      >
        Switch to Database
      </button>
    </div>
  );
}

export default App;
