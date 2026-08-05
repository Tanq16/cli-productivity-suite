-- Leader must be set before any plugin registers a mapping against it.
vim.g.mapleader = " "
vim.g.maplocalleader = " "

-- nvim-tree refuses to own the directory view unless netrw is disabled before it loads.
vim.g.loaded_netrw = 1
vim.g.loaded_netrwPlugin = 1

local o = vim.o

o.number = true
o.signcolumn = "yes"
o.cursorline = true
o.mouse = "a"
o.undofile = true
o.swapfile = false
o.ignorecase = true
o.smartcase = true
o.splitbelow = true
o.splitright = true
o.scrolloff = 8
o.sidescrolloff = 8
o.laststatus = 3
o.showmode = false
-- noinsert highlights the top candidate without writing it; noselect would leave nothing for <CR> to accept.
o.completeopt = "menu,menuone,noinsert,popup,fuzzy"
o.updatetime = 250
o.timeoutlen = 400
o.expandtab = true
o.shiftwidth = 2
o.tabstop = 2
o.softtabstop = 2
o.smartindent = true
o.wrap = false
o.confirm = true
-- Adds the arrow keys to the stock b,s so they cross line boundaries: <> in normal mode, [] in insert.
o.whichwrap = "b,s,<,>,[,]"

o.list = true
vim.opt.listchars = { tab = "│ ", leadmultispace = "│ ", trail = "·", nbsp = "␣" }
vim.opt.fillchars = { eob = " " }

-- Off means Nvim emits cterm attributes, so indices 0-15 resolve against the terminal's own palette and follow its theme.
o.termguicolors = false
-- This scheme leaves Normal undefined, which keeps the background transparent.
vim.cmd.colorscheme("vim")

local hl = {
  Comment = { ctermfg = 8, italic = true },
  Constant = { ctermfg = 14 },
  String = { ctermfg = 2 },
  Character = { ctermfg = 2 },
  Number = { ctermfg = 13 },
  Boolean = { ctermfg = 13, bold = true },
  Float = { ctermfg = 13 },
  Identifier = { ctermfg = 15 },
  Function = { ctermfg = 12 },
  Statement = { ctermfg = 5 },
  Conditional = { ctermfg = 5 },
  Repeat = { ctermfg = 5 },
  Label = { ctermfg = 5 },
  Operator = { ctermfg = 6 },
  Keyword = { ctermfg = 5 },
  Exception = { ctermfg = 9 },
  PreProc = { ctermfg = 13 },
  Include = { ctermfg = 13, bold = true },
  Define = { ctermfg = 13 },
  Macro = { ctermfg = 13 },
  Type = { ctermfg = 11 },
  StorageClass = { ctermfg = 11 },
  Structure = { ctermfg = 11, bold = true },
  Typedef = { ctermfg = 11 },
  Special = { ctermfg = 6 },
  SpecialKey = { ctermfg = 8 },
  Delimiter = { ctermfg = 7 },
  Todo = { ctermfg = 11, bold = true },
  Error = { ctermfg = 9, bold = true },
  Underlined = { ctermfg = 12, underline = true },

  LineNr = { ctermfg = 8 },
  CursorLineNr = { ctermfg = 11, ctermbg = 0, bold = true },
  CursorLine = { ctermbg = 0 },
  CursorLineSign = { ctermbg = 0 },
  CursorLineFold = { ctermbg = 0 },
  -- Vim's default fills these with grey, drawing a solid bar down the gutter of every window.
  SignColumn = {},
  FoldColumn = { ctermfg = 8 },
  Whitespace = { ctermfg = 8 },
  NonText = { ctermfg = 8 },
  Conceal = { ctermfg = 8 },
  Visual = { reverse = true },
  Search = { ctermfg = 0, ctermbg = 11 },
  IncSearch = { ctermfg = 0, ctermbg = 9 },
  MatchParen = { ctermfg = 14, bold = true, underline = true },
  Pmenu = { ctermbg = 8 },
  PmenuSel = { ctermbg = 4, ctermfg = 15, bold = true },
  PmenuSbar = { ctermbg = 8 },
  PmenuThumb = { ctermbg = 7 },
  WinSeparator = { ctermfg = 8 },
  Folded = { ctermfg = 8, italic = true },
  ErrorMsg = { ctermfg = 9, bold = true },
  TabLine = { ctermfg = 8 },
  TabLineFill = {},
  TabLineSel = { ctermfg = 4, bold = true },

  -- Vim's default is cterm=reverse, which renders the whole bar as a solid slab.
  StatusLine = { ctermfg = 7 },
  StatusLineNC = { ctermfg = 8 },
  StatusLineTerm = { ctermfg = 7 },
  StatusLineTermNC = { ctermfg = 8 },
  StlGit = { ctermfg = 5 },
  StlErr = { ctermfg = 9 },
  StlWarn = { ctermfg = 11 },

  DiagnosticError = { ctermfg = 9 },
  DiagnosticWarn = { ctermfg = 11 },
  DiagnosticInfo = { ctermfg = 12 },
  DiagnosticHint = { ctermfg = 14 },
  DiagnosticUnderlineError = { ctermfg = 9, undercurl = true },
  DiagnosticUnderlineWarn = { ctermfg = 11, undercurl = true },

  DiffAdd = { ctermfg = 2 },
  DiffChange = { ctermfg = 3 },
  DiffDelete = { ctermfg = 1 },
  DiffText = { ctermfg = 11, bold = true },
  GitSignsAdd = { ctermfg = 2 },
  GitSignsChange = { ctermfg = 3 },
  GitSignsDelete = { ctermfg = 1 },

  ["@variable"] = { ctermfg = 15 },
  ["@variable.builtin"] = { ctermfg = 9, italic = true },
  ["@variable.parameter"] = { ctermfg = 7, italic = true },
  ["@variable.member"] = { ctermfg = 14 },
  ["@function.builtin"] = { ctermfg = 12, bold = true },
  ["@function.call"] = { ctermfg = 12 },
  ["@function.method"] = { ctermfg = 12 },
  ["@constructor"] = { ctermfg = 11, bold = true },
  ["@type.builtin"] = { ctermfg = 3 },
  ["@keyword.import"] = { ctermfg = 13, bold = true },
  ["@keyword.return"] = { ctermfg = 5, bold = true },
  ["@punctuation.bracket"] = { ctermfg = 7 },
  ["@punctuation.delimiter"] = { ctermfg = 7 },
  ["@comment.documentation"] = { ctermfg = 8, italic = true },
  ["@markup.heading"] = { ctermfg = 12, bold = true },
  ["@markup.link"] = { ctermfg = 14, underline = true },
  ["@markup.raw"] = { ctermfg = 2 },
  ["@diff.plus"] = { ctermfg = 2 },
  ["@diff.minus"] = { ctermfg = 1 },
}

-- A cap is the pill colour drawn as a foreground, so every pill colour needs a background group and a matching foreground one.
for color = 1, 6 do
  hl["StlPill" .. color] = { ctermfg = 0, ctermbg = color, bold = true }
  hl["StlCap" .. color] = { ctermfg = color }
end

local function apply_highlights()
  for group, spec in pairs(hl) do
    vim.api.nvim_set_hl(0, group, spec)
  end
end

vim.api.nvim_create_autocmd("ColorScheme", { callback = apply_highlights })
apply_highlights()

vim.pack.add({
  { src = "https://github.com/ibhagwan/fzf-lua" },
  { src = "https://github.com/nvim-tree/nvim-tree.lua" },
  { src = "https://github.com/lewis6991/gitsigns.nvim" },
  { src = "https://github.com/windwp/nvim-autopairs" },
  -- The "main" rewrite has a different API than "master"; pin it so an upstream default-branch change cannot swap it silently.
  { src = "https://github.com/nvim-treesitter/nvim-treesitter", version = "main" },
})

local ts_langs = {
  "bash",
  "go",
  "gomod",
  "javascript",
  "json",
  "lua",
  "markdown",
  "markdown_inline",
  "python",
  "tsx",
  "typescript",
  "yaml",
}

-- The "main" branch shells out to the tree-sitter CLI to compile each parser.
if vim.fn.executable("tree-sitter") == 1 then
  local missing = vim.tbl_filter(function(lang)
    return #vim.api.nvim_get_runtime_file("parser/" .. lang .. ".so", false) == 0
  end, ts_langs)
  if #missing > 0 then
    require("nvim-treesitter").install(missing)
  end
end

vim.treesitter.language.register("bash", "zsh")

-- Buffers without a parser fall back to the bundled regex syntax files, so failure here is not an error.
vim.api.nvim_create_autocmd("FileType", {
  callback = function(args)
    pcall(vim.treesitter.start, args.buf)
  end,
})

require("nvim-tree").setup({
  view = { width = 34 },
  renderer = { group_empty = true, indent_markers = { enable = true } },
  filters = { dotfiles = false, custom = { "^\\.git$" } },
  actions = { open_file = { quit_on_open = true } },
})

require("gitsigns").setup({
  on_attach = function(bufnr)
    local gs = require("gitsigns")
    local function map(lhs, rhs, desc)
      vim.keymap.set("n", lhs, rhs, { buffer = bufnr, desc = desc })
    end
    map("]c", function() gs.nav_hunk("next") end, "next hunk")
    map("[c", function() gs.nav_hunk("prev") end, "prev hunk")
    map("<leader>hp", gs.preview_hunk, "preview hunk")
    map("<leader>hr", gs.reset_hunk, "reset hunk")
    map("<leader>hs", gs.stage_hunk, "stage hunk")
    map("<leader>hb", function() gs.blame_line({ full = true }) end, "blame line")
    map("<leader>hd", gs.diffthis, "diff this file")
  end,
})

require("nvim-autopairs").setup({ check_ts = true })

local fzf = require("fzf-lua")
fzf.setup({})

local servers = {
  gopls = {
    cmd = { "gopls" },
    filetypes = { "go", "gomod", "gowork", "gotmpl" },
    root_markers = { "go.work", "go.mod", ".git" },
  },
  ruff = {
    cmd = { "ruff", "server" },
    filetypes = { "python" },
    root_markers = { "pyproject.toml", "ruff.toml", ".ruff.toml", ".git" },
  },
  ts_ls = {
    cmd = { "typescript-language-server", "--stdio" },
    filetypes = { "javascript", "javascriptreact", "typescript", "typescriptreact" },
    root_markers = { "tsconfig.json", "jsconfig.json", "package.json", ".git" },
  },
}

for name, cfg in pairs(servers) do
  if vim.fn.executable(cfg.cmd[1]) == 1 then
    vim.lsp.config(name, cfg)
    vim.lsp.enable(name)
  end
end

vim.api.nvim_create_autocmd("LspAttach", {
  callback = function(args)
    local client = vim.lsp.get_client_by_id(args.data.client_id)
    if client and client:supports_method("textDocument/completion") then
      vim.lsp.completion.enable(true, client.id, args.buf, { autotrigger = true })
    end
  end,
})

vim.diagnostic.config({
  virtual_text = { spacing = 2, prefix = "●" },
  severity_sort = true,
  underline = true,
  signs = {
    text = {
      [vim.diagnostic.severity.ERROR] = "E",
      [vim.diagnostic.severity.WARN] = "W",
      [vim.diagnostic.severity.INFO] = "I",
      [vim.diagnostic.severity.HINT] = "H",
    },
  },
})

local modes = {
  n = { "NORMAL", 4 }, i = { "INSERT", 2 }, v = { "VISUAL", 5 }, V = { "V-LINE", 5 },
  ["\22"] = { "V-BLOCK", 5 }, s = { "SELECT", 5 }, S = { "S-LINE", 5 }, ["\19"] = { "S-BLOCK", 5 },
  R = { "REPLACE", 1 }, c = { "COMMAND", 3 }, r = { "PROMPT", 3 },
  ["!"] = { "SHELL", 6 }, t = { "TERMINAL", 6 },
}

-- The %{% %} wrapper re-parses this result, so the returned string may use % items but interpolated text must escape %.
local function pill(color, text)
  local cap = "%#StlCap" .. color .. "#"
  return cap .. "" .. "%#StlPill" .. color .. "#" .. text .. cap .. "" .. "%#StatusLine#"
end

function _G.CpsStatusline()
  local mode = modes[vim.api.nvim_get_mode().mode:sub(1, 1)] or { "?", 4 }
  local parts = { pill(mode[2], mode[1]) }

  local head = vim.b.gitsigns_head
  if head and head ~= "" then
    parts[#parts + 1] = "%#StlGit# " .. head:gsub("%%", "%%%%") .. "%#StatusLine#"
  end
  local label
  if vim.bo.filetype == "NvimTree" then
    label = " explorer"
  elseif vim.bo.buftype == "terminal" then
    label = " terminal"
  end
  parts[#parts + 1] = label or " %f%m%r"
  parts[#parts + 1] = "%="

  local counts = vim.diagnostic.count(0)
  local errors = counts[vim.diagnostic.severity.ERROR]
  local warnings = counts[vim.diagnostic.severity.WARN]
  if errors then
    parts[#parts + 1] = "%#StlErr#E" .. errors .. " %#StatusLine#"
  end
  if warnings then
    parts[#parts + 1] = "%#StlWarn#W" .. warnings .. " %#StatusLine#"
  end
  if not label and vim.bo.filetype ~= "" then
    parts[#parts + 1] = pill(4, vim.bo.filetype) .. " "
  end
  parts[#parts + 1] = pill(2, "%l:%c") .. " %P "

  return table.concat(parts)
end

vim.o.statusline = "%{%v:lua.CpsStatusline()%}"

local map = vim.keymap.set

map("n", "<Esc>", "<cmd>nohlsearch<cr>", { desc = "clear search highlight" })
map("n", "<C-s>", "<cmd>write<cr>", { desc = "save file" })

map("i", "<C-a>", "<Home>", { desc = "start of line" })
map("i", "<C-e>", "<End>", { desc = "end of line" })

map("n", "<C-h>", "<C-w>h", { desc = "window left" })
map("n", "<C-j>", "<C-w>j", { desc = "window down" })
map("n", "<C-k>", "<C-w>k", { desc = "window up" })
map("n", "<C-l>", "<C-w>l", { desc = "window right" })
map("n", "<leader>n", "<cmd>set nu!<cr>", { desc = "toggle line numbers" })
map("n", "<leader>rn", "<cmd>set rnu!<cr>", { desc = "toggle relative numbers" })
map("n", "<leader>x", "<cmd>bdelete<cr>", { desc = "close buffer" })

map("n", "<C-n>", "<cmd>NvimTreeToggle<cr>", { desc = "toggle file tree" })
map("n", "<leader>e", "<cmd>NvimTreeFocus<cr>", { desc = "focus file tree" })

map("n", "<leader>ff", fzf.files, { desc = "find files" })
map("n", "<leader>fw", fzf.live_grep, { desc = "live grep" })
map("n", "<leader>fa", function()
  fzf.files({ fd_opts = "--color=never --type f --hidden --no-ignore --exclude .git" })
end, { desc = "find all files (hidden + ignored)" })
map("n", "<leader>fb", fzf.buffers, { desc = "find buffers" })

map({ "n", "v" }, "<leader>fm", function()
  vim.lsp.buf.format({ async = true })
end, { desc = "format buffer" })
map("n", "<leader>ds", vim.diagnostic.setloclist, { desc = "diagnostics to loclist" })
