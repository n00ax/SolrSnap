package backup

import (
	"fmt"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"os"
	"solr-snapshot-service/backup/blob"
	"solr-snapshot-service/backup/repo"
	"solr-snapshot-service/backup/solr"
	"solr-snapshot-service/backup/zookeeper"
	"strconv"
	"strings"
)

const defaultBlobPerm = 0777
const anonBlobUser = ""
const anonBlobDescription = ""

var Commands = [][2]string{
	{"LoadRepoFromZK", "Loads local repository from Zookeeper [zkConnectStr] [repoDestPath] [user] [description]"},
	{"LoadRemoteRepoFromZK", "Loads (or updates) remote repository with current ZK  [zkConnectStr] [gitRemote] [gitUsername] [gitPassword]"},
	{"LoadBlobFromZK", "Loads blob file from Zookeeper [zkConnectStr] [blobDestPath] [user] [description] [solrZkPath]"},
	{"LoadZKFromBlob", "Loads Zookeeper from blob file [zkConnectStr] [blobSrcPath]"},
	{"CustomizeBlob", "Customizes blob file with local configuration [blobSrcPath] [blobOutPath] [newNodeName] [newNodePort]"},
	{"CreateCores", "Creates cores from blob file [blobSrcPath] [solrDataPath]"},
	{"PopulateCores", "Populates cores with assigned blob data [blobSrcPath] [solrDataPath]]"},
	{"BlobInfo", "Shows blob information"},
	{"Automated", "Runs automated tasks, enter Automated for task information"},
}

func LoadRepoFromZK(connectStr string, repoDestPath string, generateUser string, description string) {
	img := getZKBlob(connectStr, generateUser, description, "/")
	err := repo.Create(repoDestPath, img.ZKRoot)
	if err != nil {
		log.Error("Repo load failed error ", err.Error())
		os.Exit(-1)
	}
}

func LoadBlobFromZK(connectStr string, blobDestPath string, generateUser string, description string, solrZkPath string) {
	log.Info("Grabbing objects from Zookeeper")
	img := getZKBlob(connectStr, generateUser, description, solrZkPath)
	log.Info("Saving blob from Zookeeper")
	err := img.Save(blobDestPath, defaultBlobPerm)
	if err != nil {
		log.Error("Error saving blob image error ", err.Error())
		os.Exit(-1)
	}
}

func LoadZKFromBlob(connectStr string, blobSrcPath string) {
	err := zookeeper.Restore(getBlob(blobSrcPath).ZKRoot, strings.Split(connectStr, ","), blob.FusionRestExport)
	if err != nil {
		log.Error("Error saving blob image error ", err.Error())
		os.Exit(-1)
	}
}

func LoadRemoteRepoFromZK(connectStr string, gitRemote string, gitUsername string, gitPassword string) {
	img := getZKBlob(connectStr, anonBlobUser, anonBlobDescription, "/")
	err := repo.Update(gitRemote, gitUsername, gitPassword, img.ZKRoot)
	if err != nil {
		log.Error("Error applying repo updates error ", err.Error())
		os.Exit(-1)
	}
}

func BlobInfo(blobSrcPath string) {
	img := getBlob(blobSrcPath)
	property := color.New(color.FgWhite, color.Bold)
	color.Green("=== Image Header ===")
	_, _ = property.Printf("SHA256Checksum: \t%s\n", img.Header.SHA256Checksum)
	_, _ = property.Printf("Source: \t\t%s\n", img.Header.Source)
	_, _ = property.Printf("TimeGenerated: \t\t%s\n", img.Header.TimeGenerated)
	_, _ = property.Printf("Description: \t\t%s\n", img.Header.Description)
	_, _ = property.Printf("SourceApplication: \t%s\n", img.Header.SourceApplication)
	_, _ = property.Printf("Version: \t\t%s\n", img.Header.Version)
	_, _ = property.Printf("GenerateUser: \t\t%s\n", img.Header.GenerateUser)
	_, _ = property.Printf("Magic: \t\t\t%s\n", img.Header.Magic)
	_, _ = property.Printf("SolrZKPath: \t\t\t%s\n", img.Header.SolrZKPath)
	color.Green("=== ZK Storage ===")
	_, _ = property.Printf("Total Objects: \t\t%d\n", img.ZKRoot.GetObjCount())
	validation := img.ValidateImage()
	if validation != nil {
		color.Red("Image validation failed error ", validation.Error())
		os.Exit(-1)
	}
	_, _ = property.Printf("Validation Status: ")
	color.Green("\tGood!")
	color.Green("=== Index Storage ===")
	_, _ = property.Printf("Status: \t\t")
	color.Yellow("Not Assigned")
}

func CustomizeBlob(blobSrcPath string, blobOutPath string, newNodeName string, newNodePort string) {
	img := getBlob(blobSrcPath)
	port, err := strconv.Atoi(newNodePort)
	if err != nil {
		log.Error("Failed to convert Port ", port, "to int error ", err.Error())
		os.Exit(-1)
	}
	err = solr.CustomizeBlob(img, newNodeName, port)
	if err != nil {
		log.Error("Failed to customize image error ", err.Error())
		os.Exit(-1)
	}
	err = img.Save(blobOutPath, defaultBlobPerm)
	if err != nil {
		log.Error("Failed to save image error ", err.Error())
		os.Exit(-1)
	}
}

func CreateCores(blobSrcPath string, solrDataPath string) {
	img := getBlob(blobSrcPath)
	err := solr.CreateCores(img, solrDataPath)
	if err != nil {
		log.Error("Failed to create cores error ", err.Error())
		os.Exit(-1)
	}
}

func getZKBlob(connectStr string, generateUser string, description string, solrZKPath string) *blob.Image {
	servers := strings.Split(connectStr, ",")
	img, err := zookeeper.Read(servers, blob.ZKDirectExport, generateUser, description, solrZKPath)
	if err != nil {
		log.Error(err.Error())
		os.Exit(-1)
	}
	return img
}

func getBlob(filePath string) *blob.Image {
	img, err := blob.Read(filePath)
	if err != nil {
		log.Errorf(err.Error())
		os.Exit(-1)
	}
	return img
}

func Help() {
	fmt.Println("Available actions:")
	for _, mapping := range Commands {
		fmt.Printf("%s - %s\n", mapping[0], mapping[1])
	}
}
