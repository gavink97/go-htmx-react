import { createRoot } from 'react-dom/client';
import { NavBar } from '@/assets/navbar/navbar';

const headerRoot = document.getElementById('react-nav');
if (!headerRoot) {
	throw new Error('Could not find element with id react-nav');
}
const headerReactRoot = createRoot(headerRoot);
headerReactRoot.render(NavBar());
