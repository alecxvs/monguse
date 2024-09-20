package steamgamepath

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/andygrunwald/vdf"
	"github.com/go-viper/mapstructure/v2"
	"golang.org/x/sys/windows/registry"
)

type LibraryFolders struct {
	LibraryFolders map[string]LibraryFolder `mapstructure:"libraryfolders"`
}

type LibraryFolder struct {
	Path string            `mapstructure:"path"`
	Apps map[string]string `mapstructure:"apps"`
}

type AppManifest struct {
	AppState AppState `mapstructure:"appstate"`
}

type AppState struct {
	AppId      string `mapstructure:"appid"`
	InstallDir string `mapstructure:"installdir"`
}

func GetSteamPath() (string, error) {
	regk, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Valve\Steam`, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf(`unexpected error opening 'HKEY_LOCAL_MACHINE\SOFTWARE\WOW6432Node\Valve\Steam': %w`, err)
	}
	defer regk.Close()

	path, _, err := regk.GetStringValue("InstallPath")
	if err != nil {
		return "", fmt.Errorf("unexpected error reading Windows Registry to parse InstallPath: %w", err)
	}

	return path, nil
}

func GetSteamLibraries(steamPath string) (libs []LibraryFolder, err error) {
	path := path.Join(steamPath, "steamapps", "libraryfolders.vdf")
	folderVdf, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unexpected error opening file %s: %w", path, err)
	}

	dataMap, err := vdf.NewParser(folderVdf).Parse()
	if err != nil {
		return nil, fmt.Errorf("unexpected error VDF parsing %s: %w", path, err)
	}

	var result LibraryFolders
	err = mapstructure.Decode(dataMap, &result)
	if err != nil {
		return nil, fmt.Errorf("unexpected error decoding libraries mapping: %w", err)
	}

	for _, library := range result.LibraryFolders {
		libs = append(libs, library)
	}

	return
}

func GetSteamGameLibraryPath(appId int) (string, error) {
	steamPath, err := GetSteamPath()
	if err != nil {
		return "", fmt.Errorf("failed to get Steam installation path. %w", err)
	}
	libraries, err := GetSteamLibraries(steamPath)
	if err != nil {
		return "", fmt.Errorf("failed to get Steam libraries. %w", err)
	}

	var foundPaths []string
	for _, library := range libraries {
		for libraryAppId := range library.Apps {
			if libraryAppId == string(appId) {
				foundPaths = append(foundPaths, library.Path)
			}
		}
	}

	if len(foundPaths) == 0 {
		log.Printf("Warning: Did not find target game %d in installed Steam games!\n", appId)
		return "", nil
	}

	if len(foundPaths) > 1 {
		var stb strings.Builder
		stb.WriteString(fmt.Sprintf("Unexpected multiple libraries with game %d:", appId))
		for i, s := range foundPaths {
			stb.WriteString(fmt.Sprintf("  %d. '%s'", i, s))
		}
		// TODO: We should have functionality to select a path, possibly.
		//   But also, putting this message here violates its encapsulation!
		//		Two solutions:
		//			- Return multiple path strings
		//			- Custom error type for this scenario
		//   Can Steam even have the same app ID in two places!?
		//   Maybe they could add splitting game files across drives to squeeze out disk space?
		stb.WriteString("This program can't make a determination of which copy of the game to return.")
		stb.WriteString("Manually input the correct path, if possible.")
		log.Println(stb.String())
		return "", fmt.Errorf("unexpected multiple libraries with game %d: %v", appId, foundPaths)
	}

	return foundPaths[0], nil
}

func GetSteamGamePath(appId int) (string, error) {
	libraryPath, err := GetSteamGameLibraryPath(appId)
	if err != nil {
		return "", err
	}

	appManifestPath := path.Join(libraryPath, fmt.Sprintf("appmanifest_%d.acf", appId))

	appManifestVdf, err := os.Open(appManifestPath)
	if err != nil {
		return "", fmt.Errorf("unexpected error opening app manifest file (%s): %w", appManifestPath, err)
	}

	dataMap, err := vdf.NewParser(appManifestVdf).Parse()
	if err != nil {
		return "", fmt.Errorf("unexpected error VDF parsing %s: %w", appManifestPath, err)
	}

	var result AppManifest
	err = mapstructure.Decode(dataMap, &result)
	if err != nil {
		return "", fmt.Errorf("unexpected error decoding app manifest mapping: %w", err)
	}

	appPath := path.Join(libraryPath, "steamapps", "common", result.AppState.InstallDir)
	_, err = os.Stat(appPath)
	if err != nil {
		return "", fmt.Errorf("computed path '%s' does not exist! %w", appPath, err)
	}

	return appPath, nil
}
