import * as vscode from "vscode";
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node";
import * as path from "path";

let client: LanguageClient;

export function activate(context: vscode.ExtensionContext) {
  console.log("Foo Language extension activated 1.1.5");

  // Путь к LSP серверу - абсолютный путь
  const serverExecutable = "D:\\dev\\go\\foo_lang_v2\\lsp\\foo-lsp.exe";

  // Опции для запуска LSP сервера
  const serverOptions: ServerOptions = {
    run: {
      command: serverExecutable,
      transport: TransportKind.stdio,
    },
    debug: {
      command: serverExecutable,
      transport: TransportKind.stdio,
    },
  };

  // Опции языкового клиента
  const clientOptions: LanguageClientOptions = {
    // Регистрируем схему документов для активации LSP
    documentSelector: [{ scheme: "file", language: "foo" }],
    synchronize: {
      // Уведомляем сервер об изменениях в .foo файлах
      fileEvents: vscode.workspace.createFileSystemWatcher("**/*.foo"),
    },
    // Настройки инициализации
    initializationOptions: {
      settings: {
        foo: {
          maxCompletionItems: 100,
          enableDiagnostics: true,
          enableHover: true,
        },
      },
    },
  };

  // Создаем и запускаем языковой клиент
  client = new LanguageClient(
    "fooLanguageServer",
    "Foo Language Server",
    serverOptions,
    clientOptions
  );

  // Регистрируем команды расширения
  const disposable = vscode.commands.registerCommand("foo.restart", () => {
    client.restart();
    vscode.window.showInformationMessage("Foo Language Server restarted");
  });

  context.subscriptions.push(disposable);

  // Запускаем LSP клиент
  client
    .start()
    .then(() => {
      console.log("Foo Language Server started successfully");

      // Показываем уведомление о готовности
      vscode.window.showInformationMessage(
        "Foo Language Server is now active! 🚀"
      );
    })
    .catch((error) => {
      console.error("Failed to start Foo Language Server:", error);
      vscode.window.showErrorMessage(
        `Failed to start Foo Language Server: ${error.message}`
      );
    });
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  console.log("Foo Language extension deactivated");
  return client.stop();
}
