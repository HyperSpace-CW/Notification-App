package utils

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"go.uber.org/zap"
)

func GetRequestID(ctx *fiber.Ctx) string {
	headerExecID := ctx.Request().Header.Peek("Execid")
	if len(headerExecID) != 0 {
		return string(utils.CopyBytes(headerExecID))
	}
	generatedRequestID, _ := ctx.Locals("requestid").(string)
	return generatedRequestID
}

func GetPathVariableInt64(ctx *fiber.Ctx, varName string) (int64, error) {
	varStr := utils.CopyString(ctx.Params(varName))
	if varStr == "" {
		return 0, fmt.Errorf("missing %s path variable", varName)
	}
	variable, err := strconv.ParseInt(varStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("path variable %s must be digit: %w", varName, err)
	}
	return variable, nil
}

func GetQueryParameterInt64(ctx *fiber.Ctx, paramName string) (int64, error) {
	paramStr := utils.CopyString(ctx.Query(paramName))
	if paramStr == "" {
		return 0, fmt.Errorf("missing %s query parameter", paramName)
	}
	paramInt64, err := strconv.ParseInt(paramStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("query parameter %s must be digit: %w", paramName, err)
	}
	return paramInt64, nil
}

func GetQueryParameterBool(ctx *fiber.Ctx, paramName string) (bool, error) {
	paramStr := utils.CopyString(ctx.Query(paramName))
	if paramStr == "" {
		return false, fmt.Errorf("missing %s query parameter", paramName)
	}
	paramBool, err := strconv.ParseBool(paramStr)
	if err != nil {
		return false, fmt.Errorf("query parameter %s must be boolean: %w", paramName, err)
	}
	return paramBool, nil
}

func NewCtxFromFiberCtx(ctx *fiber.Ctx) context.Context {
	requestID, _ := ctx.Locals("requestid").(string)
	return context.WithValue(context.Background(), "processId", requestID) //nolint:revive,staticcheck
}

func SetContentTypeForFile(ctx *fiber.Ctx, fileName string) error {
	fileNameParts := strings.Split(fileName, ".")
	if len(fileNameParts) < 2 {
		return errors.New("file does not have an extension")
	}

	switch fileNameParts[len(fileNameParts)-1] {
	case "zip":
		ctx.Response().Header.Set(fiber.HeaderContentType, "application/zip")
	case "pdf":
		ctx.Response().Header.Set(fiber.HeaderContentType, "application/pdf")
	case "txt":
		ctx.Response().Header.Set(fiber.HeaderContentType, "text/plain")
	case "csv":
		ctx.Response().Header.Set(fiber.HeaderContentType, "text/csv")
	case "xls":
		ctx.Response().Header.Set(fiber.HeaderContentType, "application/vnd.ms-excel")
	case "xlsx":
		ctx.Response().Header.Set(fiber.HeaderContentType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	default:
		return errors.New("unsupported file type")
	}

	return nil
}

func SendZip(ctx *fiber.Ctx, log *zap.Logger, zipName, srcDirPath string, filePaths ...string) error {
	var buff bytes.Buffer
	zipWriter := zip.NewWriter(&buff)

	for _, filePath := range filePaths {
		if err := func() error {
			filePathParts := strings.Split(filePath, "/")
			fileName := filePathParts[len(filePathParts)-1]

			//nolint:exhaustruct
			entryWriter, err := zipWriter.CreateHeader(&zip.FileHeader{
				Name:     fileName,
				Method:   zip.Store,
				Modified: time.Now(),
			})
			if err != nil {
				return fmt.Errorf("create zip header: %w", err)
			}

			file, err := os.Open(srcDirPath + "/" + filePath)
			if err != nil {
				return fmt.Errorf("open file: %w", err)
			}
			defer func() {
				if err := file.Close(); err != nil {
					log.Error(
						"error closing file",
						zap.Error(err),
						zap.String("file_path", srcDirPath+"/"+filePath),
					)
				}
			}()

			fileReader := bufio.NewReader(file)
			_, err = io.Copy(entryWriter, fileReader)
			if err != nil {
				return fmt.Errorf("copy data from file to zip: %w", err)
			}
			return nil
		}(); err != nil {
			return err
		}
	}

	if err := zipWriter.Close(); err != nil {
		return fmt.Errorf("close zip writer: %w", err)
	}

	if err := SetContentTypeForFile(ctx, zipName+".zip"); err != nil {
		return fmt.Errorf("error set content type for file: %w", err)
	}

	ctx.Response().Header.Set("Content-Disposition", `attachment; filename="`+zipName+`.zip"`)
	if err := ctx.Send(buff.Bytes()); err != nil {
		return fmt.Errorf("send response: %w", err)
	}

	return nil
}
