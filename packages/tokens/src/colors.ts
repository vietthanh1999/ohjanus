export const primitiveColors = {
  gray: {
    50: '#F9FAFB',
    100: '#F3F4F6',
    200: '#E5E7EB',
    300: '#D1D5DB',
    400: '#9CA3AF',
    500: '#6B7280',
    600: '#4B5563',
    700: '#374151',
    800: '#1F2937',
    900: '#111827',
    950: '#0B0F19'
  },
  blue: {
    50: '#EFF6FF',
    100: '#DBEAFE',
    200: '#BFDBFE',
    300: '#93C5FD',
    400: '#60A5FA',
    500: '#3574F0', // JetBrains accent blue
    600: '#2563EB',
    700: '#1D4ED8',
    800: '#1E40AF',
    900: '#1E3A8A',
    950: '#172554'
  },
  green: {
    50: '#F0FDF4',
    100: '#DCFCE7',
    200: '#BBF7D0',
    300: '#86EFAC',
    400: '#4ADE80',
    500: '#57D38C', // JetBrains success green
    600: '#16A34A',
    700: '#15803D',
    800: '#166534',
    900: '#14532D',
    950: '#052E16'
  },
  yellow: {
    50: '#FEFCE8',
    100: '#FEF9C3',
    200: '#FEF08A',
    300: '#FDE047',
    400: '#FACC15',
    500: '#EDA200', // JetBrains warning yellow
    600: '#CA8A04',
    700: '#A16207',
    800: '#854D0E',
    900: '#713F12',
    950: '#422006'
  },
  red: {
    50: '#FEF2F2',
    100: '#FEE2E2',
    200: '#FECACA',
    300: '#FCA5A5',
    400: '#F87171',
    500: '#E55353', // JetBrains danger red
    600: '#DC2626',
    700: '#B91C1C',
    800: '#991B1B',
    900: '#7F1D1D',
    950: '#450A0A'
  },
  cyan: {
    50: '#ECFEFF',
    100: '#CFFAFE',
    200: '#A5F3FC',
    300: '#67E8F9',
    400: '#22D3EE',
    500: '#0891B2',
    600: '#0E7490',
    700: '#155E75',
    800: '#164E63',
    900: '#083344',
    950: '#041B26'
  }
};

export const darkSemanticTokens = {
  bg: {
    canvas: '#1E1F22',
    sidebar: '#1E1F22',
    toolbar: '#2B2D30',
    tableHeader: '#25272A',
    tableRow: '#1E1F22',
    hover: '#2B2D30',
    selected: '#2E3A4E',
    activeTab: '#1F2E4A',
    inactiveTab: 'transparent',
    floating: '#2B2D30',
    executionBlock: 'rgba(31, 67, 136, 0.45)',
    windowFrame: '#24272A'
  },
  fg: {
    primary: '#DFE1E5',
    secondary: '#9DA0A8',
    muted: '#7A7E85',
    nullValue: '#6F737A',
    disabled: '#5A5D63'
  },
  border: {
    panel: 'rgba(255, 255, 255, 0.08)',
    subtle: '#2B2D30',
    default: '#393B40',
    strong: '#4E5157',
    accent: '#3574F0',
    execution: '#3569C8'
  },
  status: {
    success: '#57D38C',
    warning: '#EDA200',
    danger: '#E55353',
    info: '#56A8F5'
  }
};

export const lightSemanticTokens = {
  bg: {
    canvas: '#FFFFFF',
    sidebar: '#F7F8FA',
    toolbar: '#EBECF0',
    tableHeader: '#F2F3F5',
    tableRow: '#FFFFFF',
    hover: '#EBECF0',
    selected: '#D4E2FF',
    activeTab: '#FFFFFF',
    inactiveTab: 'transparent',
    floating: '#FFFFFF',
    executionBlock: 'rgba(53, 116, 240, 0.15)',
    windowFrame: '#E5E7EB'
  },
  fg: {
    primary: '#1F2328',
    secondary: '#4B5563',
    muted: '#6B7280',
    nullValue: '#9CA3AF',
    disabled: '#D1D5DB'
  },
  border: {
    panel: 'rgba(0, 0, 0, 0.1)',
    subtle: '#E5E7EB',
    default: '#D1D5DB',
    strong: '#9CA3AF',
    accent: '#3574F0',
    execution: '#3569C8'
  },
  status: {
    success: '#16A34A',
    warning: '#D97706',
    danger: '#DC2626',
    info: '#2563EB'
  }
};
