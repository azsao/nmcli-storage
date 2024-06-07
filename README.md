<div align="center">

# NScli

</div>

<p align="center">
<a href="#">
<img alt="Made with GO" src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white"/>
</a>
</p>

This command line interface is designed to store your network passwords locally and encripted using the popular NetworkManager client.

# ‍ ‍ ‍ 

<details>
<summary><strong>Information</strong></summary>
This program is designed to manage Wi-Fi credentials securely by providing an interface for storing, encrypting, and connecting to Wi-Fi networks using the Network Manager CLI (nmcli).

</details> 

<details>
<summary><strong>How it works?</strong></summary>

* Once the user is prompted to input an SSID and password, these credentials are encrypted using an encryption key derived from a secure password.
* The program uses AES (Advanced Encryption Standard), which is a symmetric key encryption algorithm widely regarded for its strength and efficiency. AES is used by governments, financial institutions, and organizations worldwide to protect sensitive data.
* The program generates an AES key using a password provided by the user. This key is essential for both encrypting and decrypting the data, inorder to generate this key, the program uses the scrypt key derivation function.


</details>

## Key Features

* The program allows users to input Wi-Fi SSIDs (network names) and passwords.

* It securely stores these credentials in an encrypted format to prevent unauthorized access.

* Credentials are encrypted using AES (Advanced Encryption Standard), a robust encryption method.

* The encryption key is derived from a user-provided password using the scrypt key derivation function, ensuring that only users with the correct password can decrypt the stored credentials.

* Users can select a previously stored SSID to connect to the corresponding Wi-Fi network.

* The program is operated through a simple command-line interface, making it easy to use even for those with limited technical knowledge.



## Installation

> [!TIP]
> Functionality is only met once dependencies are met, NetworkManager and GO + addons must be installed prior to running this client.

### Option 1: Manually

- Clone the repository

```bash
git clone https://github.com/azsao/nmcli-storage.git 
```

- Enter the directory

```bash
cd path/to/cloned/repository
```

- Run the installation script (permissions may be necessary)

```bash
chmod +x installation.sh 
./installation.sh
```

### Option 2: CURL
> [!NOTE]
> WORK IN PROGRESS, COMING SOON!

## Dependencies

All the dependencies must be met inorder to ensure proper functionality with this repository:

- NetworkManager
- GOlang
- GOlang Crypto
