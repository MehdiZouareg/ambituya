# Ambituya

Ambituya is an application that allows you to control your Tuya devices. This README will guide you through the necessary steps to set up and use Ambituya.

## Current TODO 
- [x] Create an Ambilight effect
- [ ] Create a Settings UI so that we don't need to update config.yaml anymore
- [ ] Add a Sound effect that will change color based on PC sound 
- [ ] Handle Tuya V2 Instruction Set

## Prerequisites

Before getting started, you need to have a Tuya IoT Developer account. If you don't have one yet, follow the instructions below.

## Step 1: Create a Tuya IoT Developer Account

1. Go to the [Tuya IoT Developer](https://iot.tuya.com/) website.
2. Click on **Sign Up** to create an account. You can sign up with your email address or a third-party account.
3. Follow the instructions to complete the registration process.

## Step 2: Create a Tuya Project

1. Log in to your Tuya IoT Developer account.
2. In the dashboard, click on **Projects** in the left menu.
3. Click on **Create Project**.
4. Fill in the project information:
   - **Project Name**: Give your project a name (e.g., Ambituya).
   - **Industry**: Select the appropriate industry for your project.
   - **Description**: Add a description of your project (optional).
5. Click **Next** to proceed to the next step.
6. Select the types of devices you wish to manage (e.g., Lights, Sockets).
7. Click **Create** to finalize the creation of your project.

## Step 3: Obtain Your AccessID and AccessKey

1. In the dashboard, go to **Projects** and click on the project you just created.
2. In the project details, you will see your **AccessID** and **AccessKey**.
3. Make a note of this information, as you will need it for the application configuration.

## Step 4: Configure the `config.yaml` File

1. In the Ambituya project directory, locate the `config/config.yaml` file.
2. Open `config.yaml` with a text editor.
3. Fill in the file with your information as follows:

```yaml
# Tuya Configuration
accessID: "YOUR_ACCESS_ID"
accessKey: "YOUR_ACCESS_KEY"
appName: "YOUR_APP_NAME"
debugMode: false

# Devices
tuya:
  devices:
    - name: "SOME_DEVICE_NAME"
      id: "SOME_ID"
    - name: "SOME_ANOTHER_DEVICE"
      id: "SOME_ID"

# Ambilight Configuration
ambilight:
  refreshRate: 5000 # In milliseconds
```

### Side notes
This project handles only V1 Instruction sets for the time being and only the Ambilight effect. Please subscribe to this repository to be notified on new updates !

