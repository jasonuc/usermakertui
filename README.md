# **usermakertui** 🚀

**A dynamic, real-time form-based TUI built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)!**  
Showcases how to create responsive forms in the terminal, giving users instant feedback as they type.

## **Features**

- 🖥️ **Interactive Real-Time Form** - Provide immediate, real-time feedback to users as they enter data.
- 🔒 **Customizable Input Fields** - Includes secure password input, and can be easily extended with more fields.
- 🎨 **Responsive Design** - Built using [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss) for modern, sleek UIs.
- ⚙️ **Flexible Integration** - Easy to adapt and integrate into your own TUI projects.

## **Getting Started**

### **Installation & Usage**

`usermakertui` demonstrates how to build a real-time form in the terminal. It comes with a sample application that shows a user creation form, but you can modify it for other use cases.

1. Clone and build:

    ```bash
    git clone https://github.com/jasonuc/usermakertui.git
    cd usermakertui
    go build -o usermaker
    ```

2. Run the example form:

    ```bash
    ./usermaker
    ```

### **How It Works**

`usermakertui` leverages the power of [Bubble Tea](https://github.com/charmbracelet/bubbletea) to create a form that **validates inputs in real-time**. Users get immediate feedback as they type, with errors highlighted and suggestions displayed dynamically. The example provided shows how you can create a form for user input (e.g., email, password), but the concept can be extended to any kind of terminal-based form.

### **Demo**

![Demo](demo.gif)

### **Customization**

Want to adapt the real-time form for your own needs? Here’s how you can customize it:

- **Add or Modify Input Fields**: Add new `textinput.Model` components for different types of data.
- **Change Styles**: Easily tweak the styling using [Lip Gloss](https://github.com/charmbracelet/lipgloss) to match your brand or aesthetic.
- **Use for Different Applications**: Integrate this form concept into other Bubble Tea-based TUIs, like setup wizards, data entry tools, or interactive scripts.

### **Why Real-Time Feedback in TUIs?**

The real power of `usermakertui` is showing how **real-time feedback** can transform the user experience in terminal applications. By giving users immediate validation, error handling, and visual cues, you can make your TUIs as intuitive and user-friendly as modern graphical apps.

### **Open Source & Contributions**

`usermakertui` is open source and free to use. Feel free to explore, adapt, and share your modifications. Contributions are welcome — fork the repository, make your improvements, and open a pull request.
