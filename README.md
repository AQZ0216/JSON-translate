# JSON-translate

## GCP 服務帳戶金鑰申請
* https://cloud.google.com/translate/docs/quickstart-client-libraries
* 照著"事前準備"步驟申請
* ` export GOOGLE_APPLICATION_CREDENTIALS="[PATH]" `
* Ex: ` export GOOGLE_APPLICATION_CREDENTIALS="/home/user/Downloads/service-account-file.json" `

## translateJSON function
* func translateJSON(source []byte, target []byte, targetLanguage string) ([]byte, error)
  > * source : the source JSON file
  > * target : the target JSON file
  > * targetLanguage : the target language you want to translate
  >   * (ISO -639-1 代碼 : https://cloud.google.com/translate/docs/languages)
  > * return : the translated JSON file
  
# reference
* https://cloud.google.com/translate/docs/
