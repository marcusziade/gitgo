# GitGo: GitHub Issue Importer

GitGo is a simple yet powerful tool designed to automate the process of importing tasks from a CSV file into GitHub issues. It streamlines the workflow for developers and teams looking to migrate their task management into GitHub without the manual hassle.

## Features

-   **Simple CSV Import**: Easily import tasks from a CSV file directly into GitHub issues.
-   **Customizable Import Options**: Tailor the import process to fit your project's needs.
-   **Secure Authentication**: Utilizes GitHub personal access tokens for secure API access.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

-   Go (version 1.15 or higher)
-   GitHub account and a personal access token

### Installation

1. **Clone the Repository**

    ```sh
    git clone https://github.com/marcusziade/gitgo.git
    cd gitgo
    ```

2. **Set Up Your GitHub Personal Access Token**

    Add your GitHub personal access token to your environment variables:

    - For macOS/Linux:

        Add the following line to your `~/.bash_profile` or `~/.bashrc`:

        ```sh
        export GITHUB_TOKEN="your_personal_access_token"
        ```

        Reload the profile:

        ```sh
        source ~/.bash_profile
        ```

    - For Windows:

        Add the `GITHUB_TOKEN` variable to your environment variables through the System Properties.

3. **Build the Project**

    Navigate to the project directory and build the project:

    ```sh
    go build
    ```

### Usage

1. **Prepare Your CSV File**

    Ensure your CSV file is in the format: `title, description` without headers. Save it as `tasks.csv`.

2. **Run the Importer**

    Execute the compiled binary to start importing tasks:

    ```sh
    ./gitgo
    ```

## Configuration

You can customize the import process by modifying the `import_issues.go` script. Currently, it's set up to read `tasks.csv` from the same directory. Make sure to adjust the repository owner and name within the script to match your GitHub repository details.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE) file for details.
